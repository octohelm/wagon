package task

import (
	"context"
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
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
		contents := c.Directory(dagger.DirectoryOpts{
			ID: cp.Contents.DirectoryID(),
		})

		if source := cp.Source; source != "/" {
			// When file exists
			if f, err := contents.File(source).Sync(ctx); err == nil {
				out := c.
					Directory(dagger.DirectoryOpts{
						ID: cp.Input.DirectoryID(),
					}).
					WithFile(cp.Dest, f)

				return cp.Output.SetDirectoryIDBy(ctx, out)
			}

			contents = contents.Directory(source)
		}

		ct := c.
			Directory(dagger.DirectoryOpts{
				ID: cp.Input.DirectoryID(),
			}).
			WithDirectory(cp.Dest, contents, dagger.DirectoryWithDirectoryOpts{
				Include: cp.Include,
				Exclude: cp.Exclude,
			})

		return cp.Output.SetDirectoryIDBy(ctx, ct)
	})
}
