package daggerutil

import (
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

	"github.com/dagger/dagger/core/pipeline"
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

var DefaultRunnerHost = "docker-image://ghcr.io/dagger/engine:v0.5.1"

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
			shouldIgnore := false

			for _, v := range ev.Vertexes {
				if v.Started == nil && v.Completed == nil {
					shouldIgnore = true
					continue
				}

				customName := &pipeline.CustomName{}
				if json.Unmarshal([]byte(v.Name), customName) == nil {
					v.Name = customName.Name
				}

				if v.ProgressGroup != nil {
					pp := pipeline.Path{}
					// v.ProgressGroup.Name should always use Pipeline name if exists
					if json.Unmarshal([]byte(v.ProgressGroup.Id), &pp) == nil {
						for _, p := range pp {
							if strings.HasPrefix(p.Name, PipelinePrefix) {
								v.ProgressGroup.Name = p.Name
								break
							}
						}
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

			if shouldIgnore {
				continue
			}

			cleanCh <- ev
		}

		return nil
	})

	eg.Go(func() error {
		logOutput := forwardTo(logr.FromContext(ctx), startOpts.SessionToken)

		warn, err := progressui.DisplaySolveStatus(context.Background(), nil, logOutput, cleanCh)
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

func forwardTo(l logr.Logger, scope string) io.Writer {
	return &printer{
		l:     l,
		scope: PipelinePrefix + scope,
	}
}

type printer struct {
	pipelineNames sync.Map
	l             logr.Logger
	scope         string
}

func (o *printer) Write(p []byte) (n int, err error) {
	line := string(bytes.TrimSpace(p))
	idAndMsg := strings.SplitN(line, " ", 2)
	id := idAndMsg[0]

	if len(idAndMsg) == 2 {
		msg := idAndMsg[1]

		if strings.HasPrefix(msg, o.scope) {
			parts := strings.Split(msg[len(o.scope):], " / ")

			o.pipelineNames.Store(id, parts[0])

			if len(parts) == 2 {
				msg = strings.TrimSpace(parts[1])
			} else {
				msg = ""
			}
		} else {
			if strings.Contains(msg, "export") {
				o.pipelineNames.Store(id, "Exporting")
			}
		}

		if msg != "" {
			if name, ok := o.pipelineNames.Load(id); ok {
				o.l.WithValues("name", name).Info(fmt.Sprintf("%s %s", id, msg))
			}
		}
	}

	return len(p), nil
}
