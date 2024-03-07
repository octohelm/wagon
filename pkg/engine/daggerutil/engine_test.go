package daggerutil

import (
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"os"
	"testing"
	"time"

	"dagger.io/dagger"
	"github.com/go-courier/logr"
	"github.com/go-courier/logr/slog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func TestDebugDagger(t *testing.T) {
	runnerHost := os.Getenv("BUILDKIT_HOST")

	if runnerHost == "" {
		t.Skip()
	}

	t.Run(fmt.Sprintf("With"), func(t *testing.T) {
		_ = os.Chdir("../../..")

		ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))

		engineConn, release, _ := ConnectEngine(
			ctx,
			WithRunnerHost(""),
		)
		defer release()

		client, err := dagger.Connect(ctx, dagger.WithConn(engineConn))
		if err != nil {
			panic(errors.Wrapf(err, "connect dagger failed"))
		}
		defer client.Close()

		cc := client.Pipeline("$$pipeline")

		c := cc.Container().
			From("busybox").
			WithEnvVariable("DATE", time.Now().String()).
			WithExec([]string{"sh", "-c", "mkdir -p /dist"}).
			WithExec([]string{"sh", "-c", "echo ${DATE} > /dist/txt"})

		dir, err := c.Rootfs().Sync(ctx)
		if err != nil {
			panic(err)
		}

		id, err := dir.ID(ctx)
		if err != nil {
			panic(gqlerror.List{})
		}

		c2 := cc.Container().WithRootfs(cc.LoadDirectoryFromID(id))

		// #Copy fs to local
		_, err = c2.Directory("/dist").Export(ctx, ".wagon/demo")
		if err != nil {
			panic(errors.Wrapf(err, "export to client failed"))
		}
	})
}
