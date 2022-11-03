package daggerutil

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/ctxutil"
	"golang.org/x/net/context"
)

var ClientContext = ctxutil.New[*dagger.Client]()

func Do(ctx context.Context, do func(c *dagger.Client) error) (e error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				e = x
			default:
				e = errors.Errorf("%v", x)
			}
		}
	}()
	return do(ClientContext.From(ctx))
}

func ConnectDo(ctx context.Context, do func(ctx context.Context) error) error {
	runnerHost := os.Getenv("BUILDKIT_HOST")

	if buildkitArch := os.Getenv("BUILDKIT_ARCH"); buildkitArch != "" {
		runnerHost = os.Getenv("BUILDKIT_HOST_" + strings.ToUpper(buildkitArch))
	} else {
		runnerHost = os.Getenv("BUILDKIT_HOST")
	}

	if err := StartEngineOnBackground(ctx, WithRunnerHost(runnerHost)); err != nil {
		return err
	}

	c, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	newCtx, cancel := context.WithCancel(ctx)

	go func() {
		<-signalCh
		cancel()
	}()

	return do(ClientContext.Inject(newCtx, c))
}
