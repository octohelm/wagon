package daggerutil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dagger/dagger/core"
	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/router"
	"github.com/go-courier/logr"
	"github.com/google/uuid"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

var PipelinePrefix = "__PIPELINE__"

type engineOption struct {
	RunnerHost string
}

type EngineOptionFunc = func(x *engineOption)

func WithRunnerHost(runnerHost string) EngineOptionFunc {
	return func(x *engineOption) {
		x.RunnerHost = runnerHost
	}
}

var DefaultRunnerHost = "docker-image://ghcr.io/dagger/engine:v0.3.12"

func RunnerHost() string {
	var runnerHost string
	if v, ok := os.LookupEnv("_EXPERIMENTAL_DAGGER_RUNNER_HOST"); ok {
		runnerHost = v
	} else {
		runnerHost = DefaultRunnerHost
	}
	return runnerHost
}

func StartEngineOnBackground(ctx context.Context, optFns ...EngineOptionFunc) error {
	opt := &engineOption{}
	for i := range optFns {
		optFns[i](opt)
	}

	if opt.RunnerHost == "" {
		opt.RunnerHost = RunnerHost()
	}

	startOpts := &engine.Config{
		RunnerHost:        opt.RunnerHost,
		RawBuildkitStatus: make(chan *client.SolveStatus),
	}

	startOpts.SessionToken = uuid.Must(uuid.NewRandom()).String()

	eg := &errgroup.Group{}

	cleanCh := make(chan *client.SolveStatus)
	eg.Go(func() error {
		defer close(cleanCh)

		for ev := range startOpts.RawBuildkitStatus {
			for _, v := range ev.Vertexes {
				customName := &core.CustomName{}
				if json.Unmarshal([]byte(v.Name), customName) == nil {
					v.Name = customName.Pipeline.String()
				}

				if v.ProgressGroup != nil {
					pp := core.PipelinePath{}
					// v.ProgressGroup.Name should always use Pipeline name
					if json.Unmarshal([]byte(v.ProgressGroup.Id), &pp) == nil {
						v.ProgressGroup.Name = pp[0].Name
					}
				}

				if v.Completed == nil && v.Started != nil {
					// added to statuses for logging
					if strings.HasPrefix(v.Name, "exec") {
						ev.Statuses = append(ev.Statuses, &client.VertexStatus{
							ID:        v.Name,
							Vertex:    v.Digest,
							Started:   v.Started,
							Completed: v.Completed,
						})
					}
				}
			}

			cleanCh <- ev
		}

		return nil
	})

	eg.Go(func() error {
		logOutput := forwardTo(logr.FromContext(ctx))

		warn, err := progressui.DisplaySolveStatus(context.Background(), "", nil, logOutput, cleanCh)
		for _, w := range warn {
			_, _ = fmt.Fprintf(logOutput, "=> %s\n", w.Short)
		}
		return err
	})

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	_ = os.Setenv("DAGGER_SESSION_TOKEN", startOpts.SessionToken)
	_ = os.Setenv("DAGGER_SESSION_PORT", fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port))

	srv := http.Server{
		ReadHeaderTimeout: 30 * time.Second,
	}

	eg.Go(func() error {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		<-signalCh
		if err := srv.Shutdown(ctx); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				logr.FromContext(ctx).Error(err)
			}
		}

		return nil
	})

	eg.Go(func() error {
		logr.FromContext(ctx).
			WithValues(
				"name", "Engine",
				"runner", startOpts.RunnerHost,
			).
			Info(fmt.Sprintf("engine serve on %s", l.Addr().String()))

		return engine.Start(ctx, startOpts, func(ctx context.Context, r *router.Router) error {
			srv.Handler = r
			err := srv.Serve(l)
			if err != nil && !errors.Is(err, net.ErrClosed) {
				return err
			}
			return nil
		})
	})

	go func() {
		_ = eg.Wait()
	}()

	return nil
}

func forwardTo(l logr.Logger) io.Writer {
	return &printer{
		l: l,
	}
}

type printer struct {
	l             logr.Logger
	pipelineNames sync.Map
}

func (o *printer) Write(p []byte) (n int, err error) {
	line := string(bytes.TrimSpace(p))
	idAndInfo := strings.SplitN(line, " ", 2)
	id := idAndInfo[0]

	if len(idAndInfo) == 2 {
		info := idAndInfo[1]

		if strings.HasPrefix(idAndInfo[1], PipelinePrefix) {
			parts := strings.Split(info[len(PipelinePrefix):], " / ")

			o.pipelineNames.Store(id, parts[0])

			if len(parts) == 2 {
				info = parts[1]
			} else {
				info = ""
			}
		}

		scanner := bufio.NewScanner(bytes.NewBufferString(info))
		for scanner.Scan() {
			if l := scanner.Text(); len(l) > 0 {
				if name, ok := o.pipelineNames.Load(id); ok {
					o.l.WithValues("name", name).Info(l)
				} else {
					o.l.Info(l)
				}
			}
		}
	}

	return len(p), nil
}
