package engine

import (
	"os"
	"strings"

	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"golang.org/x/net/context"
)

type Engine struct {
}

func (e *Engine) Run(ctx context.Context) error {
	runnerHost := os.Getenv("BUILDKIT_HOST")

	if buildkitArch := os.Getenv("BUILDKIT_ARCH"); buildkitArch != "" {
		runnerHost = os.Getenv("BUILDKIT_HOST_" + strings.ToUpper(buildkitArch))
	} else {
		runnerHost = os.Getenv("BUILDKIT_HOST")
	}

	return daggerutil.StartEngineOnBackground(
		ctx,
		daggerutil.WithRunnerHost(runnerHost),
	)
}
