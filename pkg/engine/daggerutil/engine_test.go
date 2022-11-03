package daggerutil

import (
	"os"
	"testing"
	"time"

	"dagger.io/dagger"
	"github.com/go-courier/logr"
	"github.com/go-courier/logr/slog"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func _TestDebugDagger(t *testing.T) {
	_ = os.Chdir("../../..")

	ctx := logr.WithLogger(context.Background(), slog.Logger(slog.Default()))

	_ = StartEngineOnBackground(
		ctx,
		WithRunnerHost(os.Getenv("BUILDKIT_HOST")),
	)

	client, err := dagger.Connect(ctx)
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

	if _, err := c.ExitCode(ctx); err != nil {
		panic(err)
	}

	// #Copy fs to local
	_, err = c.Directory("/dist").Export(ctx, ".wagon/demo")
	if err != nil {
		panic(errors.Wrapf(err, "export to client failed"))
	}
}
