package daggerutil

import (
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"

	"dagger.io/dagger"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/ctxutil"
	"github.com/octohelm/wagon/pkg/logutil"
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

	end := logutil.LogProxyContext.From(ctx).Bind(logr.FromContext(ctx))
	defer end()

	return do(ClientContext.From(ctx))
}

func ConnectDo(ctx context.Context, do func(ctx context.Context) error) error {
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
