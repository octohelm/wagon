package daggerutil

import (
	"fmt"
	"github.com/dagger/dagger/engine"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime/debug"

	"github.com/dagger/dagger/engine/client"
	"github.com/go-courier/logr"
	"github.com/vito/progrock/console"
	"golang.org/x/net/context"
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

type EngineOptionFunc = func(x *client.Params)

func WithRunnerHost(runnerHost string) EngineOptionFunc {
	return func(x *client.Params) {
		x.RunnerHost = runnerHost
	}
}

var engineVersion = func() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if dep.Path == "github.com/dagger/dagger" {
				engine.Version = dep.Version
				return dep.Version
			}
		}
	}
	return ""
}()

var DefaultRunnerHost = fmt.Sprintf("docker-image://ghcr.io/dagger/engine:%s", engineVersion)

func RunnerHost() string {
	var runnerHost string
	if v, ok := os.LookupEnv("_EXPERIMENTAL_DAGGER_RUNNER_HOST"); ok {
		runnerHost = v
	} else {
		runnerHost = DefaultRunnerHost
	}
	return runnerHost
}

func ConnectEngine(ctx context.Context, optFns ...EngineOptionFunc) (DirectConn, func() error, error) {
	params := &client.Params{}
	for i := range optFns {
		optFns[i](params)
	}

	if params.RunnerHost == "" {
		params.RunnerHost = RunnerHost()
	}

	params.ProgrockWriter = console.NewWriter(
		os.Stdout,
		console.ShowInternal(InternalDebugFromContext(ctx)),
		console.WithUI(console.DefaultUI),
	)

	params.EngineNameCallback = func(name string) {
		logr.FromContext(ctx).Info(fmt.Sprintf("Connected to engine %s", name))
	}

	engineClient, ctx, err := client.Connect(ctx, *params)
	if err != nil {
		return nil, nil, err
	}

	return EngineConn(engineClient), func() error {
		return engineClient.Close()
	}, nil
}

func EngineConn(engineClient *client.Client) DirectConn {
	return func(req *http.Request) (*http.Response, error) {
		req.SetBasicAuth(engineClient.SecretToken, "")
		resp := httptest.NewRecorder()
		engineClient.ServeHTTP(resp, req)
		return resp.Result(), nil
	}
}

type DirectConn func(*http.Request) (*http.Response, error)

func (f DirectConn) Do(r *http.Request) (*http.Response, error) {
	return f(r)
}

func (f DirectConn) Host() string {
	return ":mem:"
}

func (f DirectConn) Close() error {
	return nil
}
