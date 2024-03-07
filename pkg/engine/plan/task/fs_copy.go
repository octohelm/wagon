package task

import (
	"context"
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"strings"
)

func init() {
	core.DefaultFactory.Register(&Copy{})
}

type Copy struct {
	core.Task

	Input    core.FS  `json:"input"`
	Contents core.FS  `json:"contents"`
	Source   string   `json:"source" default:"/"`
	Dest     string   `json:"dest" default:"/"`
	Include  []string `json:"include"`
	Exclude  []string `json:"exclude"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (cp *Copy) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		contents := cp.Contents.LoadDirectory(c)

		src, err := contents.Directory(cp.Source).Sync(ctx)
		if err != nil {
			// path /dist/txt is a file, not a directory
			if !strings.Contains(err.Error(), "not a directory") {
				return err
			}
			// try copy file
			f, err := contents.File(cp.Source).Sync(ctx)
			if err != nil {
				return err
			}
			return cp.Output.SetDirectoryIDBy(ctx, cp.Input.LoadDirectory(c).WithFile(cp.Dest, f))
		}

		ct := cp.Input.LoadDirectory(c).
			WithDirectory(cp.Dest, src, dagger.DirectoryWithDirectoryOpts{
				Include: cp.Include,
				Exclude: cp.Exclude,
			})

		return cp.Output.SetDirectoryIDBy(ctx, ct)
	})
}
