package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&WriteFile{})
}

type WriteFile struct {
	core.Task

	Input    core.FS `json:"input"`
	Path     string  `json:"path"`
	Contents string  `json:"contents"`

	Permissions int `json:"permissions" default:"0o644"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (e *WriteFile) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := e.Input.Directory(c).WithNewFile(e.Path, e.Contents, dagger.DirectoryWithNewFileOpts{
			Permissions: e.Permissions,
		})
		return e.Output.SetDirectoryIDBy(ctx, dir)
	})
}
