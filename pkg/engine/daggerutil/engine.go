package daggerutil

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/router"
	"github.com/go-courier/logr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vito/progrock"
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

var DefaultRunnerHost = "docker-image://ghcr.io/dagger/engine:v0.5.3"

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

	token := uuid.Must(uuid.NewRandom()).String()

	p := printerFor(logr.FromContext(ctx), token)

	startOpts := engine.Config{
		RunnerHost:     opt.RunnerHost,
		SessionToken:   token,
		ProgrockWriter: p,
	}

	eg := &errgroup.Group{}

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

func printerFor(l logr.Logger, scope string) progrock.Writer {
	return &printer{
		l:     l,
		scope: PipelinePrefix + scope,
	}
}

type printer struct {
	l              logr.Logger
	scope          string
	vertexGroups   sync.Map
	vertexNames    sync.Map
	groupRelations sync.Map
}

func (p *printer) Close() error {
	return nil
}

func (p *printer) WriteStatus(update *progrock.StatusUpdate) error {
	for _, m := range update.Memberships {
		for _, vid := range m.Vertexes {
			if _, ok := p.vertexGroups.Load(vid); !ok {
				p.vertexGroups.Store(vid, m.Group)
			}
		}
	}

	for _, g := range update.Groups {
		if parent := g.Parent; parent != nil {
			p.groupRelations.Store(g.Id, *parent)
		}
	}

	vertexIDs := map[string]bool{}

	for _, v := range update.Vertexes {
		if v.Internal || v.Canceled {
			continue
		}

		if v.Started == nil && v.Completed == nil {
			continue
		}

		if vertexIDs[v.Id] {
			continue
		}

		vertexIDs[v.Id] = true

		l := p.loggerFor(v.Id)

		if v.Cached {
			if v.Completed != nil {
				l.WithValues("cost", v.Duration(), "state", "CACHED").Info(v.Name)
			}
		} else if v.Completed != nil {
			if v.Error != nil {
				l.WithValues("cost", v.Duration(), "state", "DONE").Info(v.Name)
			}
		} else {
			l.Info(v.Name)
			p.vertexNames.Store(v.Id, v.Name)
		}
	}

	for _, t := range update.Tasks {
		if t.Started == nil && t.Completed == nil {
			continue
		}

		if t.Completed != nil {
			name := t.Name
			if strings.HasPrefix(name, "sha256:") {
				if n, ok := p.vertexNames.Load(t.Vertex); ok {
					name = n.(string)
				}
			}

			l := p.loggerFor(t.Vertex).WithValues("cost", t.Duration())

			if t.Total > 0 {
				l.Info(fmt.Sprintf("%s %s/%s", name, FileSize(t.Current), FileSize(t.Total)))
			} else {
				l.Info(fmt.Sprintf("%s %s", name, FileSize(t.Current)))
			}
		}
	}

	for _, log := range update.Logs {
		p.loggerFor(log.Vertex).Info(string(bytes.TrimSpace(log.Data)))
	}

	return nil
}

func (p *printer) resolveGroupID(vid string) (string, bool) {
	if g, ok := p.vertexGroups.Load(vid); ok {
		groupID := g.(string)

		for {
			if strings.HasPrefix(groupID, PipelinePrefix) {
				return groupID, ok
			}
			if parentGroupID, ok := p.groupRelations.Load(groupID); ok {
				return parentGroupID.(string), ok
			} else {
				break
			}
		}
	}
	return "", false
}

func (p *printer) loggerFor(vid string) logr.Logger {
	l := p.l

	if groupID, ok := p.resolveGroupID(vid); ok {
		if strings.HasPrefix(groupID, PipelinePrefix) {
			groupID = groupID[len(p.scope):]
		}
		l = p.l.WithValues("name", strings.Split(groupID, "@")[0])
	}

	return l
}

type FileSize int64

func (f FileSize) String() string {
	b := int64(f)
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
