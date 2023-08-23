package task

import (
	"context"
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Mkdir{})
}

type Mkdir struct {
	core.Task

	Input core.FS `json:"input"`

	Path        core.StringOrSlice `json:"path"`
	Permissions int                `json:"permissions" default:"0o755"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (e *Mkdir) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := c.Directory(dagger.DirectoryOpts{
			ID: e.Input.DirectoryID(),
		})

		for _, p := range e.Path.Values {
			dir = dir.WithNewDirectory(
				p,
				dagger.DirectoryWithNewDirectoryOpts{
					Permissions: e.Permissions,
				},
			)
		}

		return e.Output.SetDirectoryIDBy(ctx, dir)
	})
}
