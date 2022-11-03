package daggerutil

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/router"
	"github.com/go-courier/logr"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type engineOption struct {
	RunnerHost   string
	SessionToken string
	Port         uint16
	LogOutput    io.Writer
}

type EngineOptionFunc = func(x *engineOption)

func WithRunnerHost(runnerHost string) EngineOptionFunc {
	return func(x *engineOption) {
		x.RunnerHost = runnerHost
	}
}

func WithLogOutput(output io.Writer) EngineOptionFunc {
	return func(x *engineOption) {
		x.LogOutput = output
	}
}

func WithPort(port uint16) EngineOptionFunc {
	return func(x *engineOption) {
		x.Port = port
	}
}

func StartEngineOnBackground(ctx context.Context, optFns ...EngineOptionFunc) error {
	log := logr.FromContext(ctx).WithValues("name", "#Engine")

	opt := &engineOption{
		SessionToken: "-",
	}
	for i := range optFns {
		optFns[i](opt)
	}

	startOpts := &engine.Config{
		RunnerHost:   opt.RunnerHost,
		SessionToken: opt.SessionToken,
		LogOutput:    opt.LogOutput,
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	_ = os.Setenv("DAGGER_SESSION_TOKEN", opt.SessionToken)
	_ = os.Setenv("DAGGER_SESSION_PORT", fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port))

	srv := http.Server{
		ReadHeaderTimeout: 30 * time.Second,
	}

	go func() {
		<-signalCh

		close(startOpts.RawBuildkitStatus)

		if err := srv.Shutdown(ctx); err != nil {
			if !errors.Is(err, net.ErrClosed) {
				logr.FromContext(ctx).Error(err)
			}
		}
	}()

	log.Info("engine serve on %s with %s", l.Addr().String(), startOpts.RunnerHost)

	return engine.Start(ctx, startOpts, func(ctx context.Context, r *router.Router) error {
		srv.Handler = r
		err := srv.Serve(l)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			return err
		}
		return nil
	})
}
