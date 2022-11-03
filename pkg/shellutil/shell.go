package shellutil

import (
	"os"
	"strings"

	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/logutil"
	"golang.org/x/net/context"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Exec(ctx context.Context, script string, options ...interp.RunnerOption) error {
	sh, err := syntax.NewParser().Parse(strings.NewReader(script), script)
	if err != nil {
		return err
	}

	logger := logr.FromContext(ctx)

	runner, err := interp.New(
		append([]interp.RunnerOption{
			interp.StdIO(
				os.Stdin,
				logutil.Forward(logger.Debug),
				logutil.Forward(logger.Debug),
			),
		}, options...)...,
	)
	if err != nil {
		return err
	}
	return runner.Run(ctx, sh)
}
