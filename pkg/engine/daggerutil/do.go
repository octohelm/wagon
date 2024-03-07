package daggerutil

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dagger.io/dagger"
	contextx "github.com/octohelm/x/context"
	"golang.org/x/net/context"
)

var ClientContext = contextx.New[*dagger.Client]()

func Do(ctx context.Context, do func(c *dagger.Client) error) (e error) {
	return do(ClientContext.From(ctx))
}

func ConnectDo(ctx context.Context, do func(ctx context.Context) error) error {
	runnerHost := os.Getenv("BUILDKIT_HOST")

	if buildkitArch := os.Getenv("BUILDKIT_ARCH"); buildkitArch != "" {
		runnerHost = os.Getenv("BUILDKIT_HOST_" + strings.ToUpper(buildkitArch))
	} else {
		runnerHost = os.Getenv("BUILDKIT_HOST")
	}

	engineConn, release, err := ConnectEngine(ctx, WithRunnerHost(runnerHost))
	if err != nil {
		return err
	}
	defer func() {
		_ = release()
	}()

	c, err := dagger.Connect(ctx, dagger.WithConn(engineConn))
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Close()
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	newCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-signalCh
		cancel()
	}()

	return do(ClientContext.Inject(newCtx, c))
}
