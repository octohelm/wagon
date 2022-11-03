package engine

import (
	"testing"

	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/spec"
	testingx "github.com/octohelm/x/testing"
	"golang.org/x/net/context"
)

func TestShip(t *testing.T) {
	root := testingx.ProjectRoot()

	c, err := New(WithWorkdir(root), WithProjectRoot("cmd/hello"))
	if err != nil {
		PrintCueErrorIfNeed(err)
	}
	testingx.Expect(t, err, testingx.Be[error](nil))

	ctx := logr.WithLogger(context.Background(), logr.StdLogger())

	err = c.Run(ctx, spec.Platform{OS: "linux", Arch: "arm64"})

	testingx.Expect(t, err, testingx.Be[error](nil))
}
