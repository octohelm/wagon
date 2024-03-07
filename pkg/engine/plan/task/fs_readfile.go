package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&ReadFile{})
}

type ReadFile struct {
	core.Task
	Input    core.FS `json:"input"`
	Path     string  `json:"path"`
	Contents string  `json:"-" wagon:"generated,name=contents"`
}

func (e *ReadFile) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := e.Input.LoadDirectory(c)

		f := dir.File(e.Path)

		contents, err := f.Contents(ctx)
		if err != nil {
			return err
		}
		e.Contents = contents
		return nil
	})
}
