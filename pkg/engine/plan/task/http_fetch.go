package task

import (
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&HTTPFetch{})
}

type HTTPFetch struct {
	core.Task

	Source string `json:"source"`
	Dest   string `json:"dest"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (input *HTTPFetch) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := c.Directory().WithFile(input.Dest, c.HTTP(input.Source))
		return input.Output.SetDirectoryIDBy(ctx, dir)
	})
}
