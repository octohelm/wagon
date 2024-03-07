package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Rm{})
}

type Rm struct {
	core.Task

	Input core.FS `json:"input"`
	Path  string  `json:"path"`

	Output core.FS `json:"-" wagon:"generated,name=output"`

	AllowWildcard bool `json:"allowWildcard" default:"true" wagon:"deprecated"`
}

func (e *Rm) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := e.Input.LoadDirectory(c)

		newDir := dir.WithoutDirectory(e.Path)

		return e.Output.SetDirectoryIDBy(ctx, newDir)
	})
}
