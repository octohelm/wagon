package engine

import (
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/logutil"
	"golang.org/x/net/context"
	"os"
	"strings"
)

type Engine struct {
	l logutil.LogProxy
}

func (e *Engine) InjectContext(ctx context.Context) context.Context {
	return logutil.LogProxyContext.Inject(ctx, &e.l)
}

func (e *Engine) Run(ctx context.Context) error {
	runnerHost := os.Getenv("BUILDKIT_HOST")

	if buildkitArch := os.Getenv("BUILDKIT_ARCH"); buildkitArch != "" {
		runnerHost = os.Getenv("BUILDKIT_HOST_" + strings.ToUpper(buildkitArch))
	} else {
		runnerHost = os.Getenv("BUILDKIT_HOST")
	}

	if runnerHost != "" {
		// pipeline here always step by step
		// so here will switch log when each step, group the logs in matched step.
		go func() {
			err := daggerutil.StartEngineOnBackground(
				ctx,
				daggerutil.WithRunnerHost(runnerHost),
				daggerutil.WithLogOutput(logutil.Forward(e.l.Info)),
			)
			if err != nil {
				panic(err)
			}
		}()
	}

	return nil
}
