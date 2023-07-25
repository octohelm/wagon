package daggerutil

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/router"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vito/progrock/console"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type contextInternalDebug struct {
}

func ContextWithInternalDebug(ctx context.Context, debug bool) context.Context {
	return context.WithValue(ctx, contextInternalDebug{}, debug)
}

func InternalDebugFromContext(ctx context.Context) bool {
	v, ok := ctx.Value(contextInternalDebug{}).(bool)
	if ok {
		return v
	}
	return false
}

type engineOption struct {
	RunnerHost string
}

type EngineOptionFunc = func(x *engineOption)

func WithRunnerHost(runnerHost string) EngineOptionFunc {
	return func(x *engineOption) {
		x.RunnerHost = runnerHost
	}
}

var DefaultRunnerHost = "docker-image://ghcr.io/dagger/engine:v0.6.4"

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

	startOpts := engine.Config{
		RunnerHost:   opt.RunnerHost,
		SessionToken: token,
		ProgrockWriter: console.NewWriter(
			os.Stdout,
			console.ShowInternal(InternalDebugFromContext(ctx)),
			console.WithUI(console.DefaultUI),
		),
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
				return err
			}
		}
		return nil
	})

	eg.Go(func() error {
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
