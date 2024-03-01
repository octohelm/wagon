package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Merge{})
}

type Merge struct {
	core.Task

	Inputs []core.FS `json:"inputs"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (e *Merge) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		d := c.Directory()

		for _, input := range e.Inputs {
			d = d.WithDirectory("/", input.Directory(c))
		}

		return e.Output.SetDirectoryIDBy(ctx, d)
	})
}
