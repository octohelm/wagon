package task

import (
	"context"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&ImageCopy{})
}

type ImageCopy struct {
	core.Task

	Input    core.Image `json:"input"`
	Contents core.FS    `json:"contents"`
	Source   string     `json:"source" default:"/"`
	Dest     string     `json:"dest" default:"."`
	Include  []string   `json:"include"`
	Exclude  []string   `json:"exclude"`

	Output core.Image `json:"-" wagon:"generated,name=output"`
}

func (cp *ImageCopy) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		contents := c.Directory(dagger.DirectoryOpts{
			ID: cp.Contents.DirectoryID(),
		})

		if source := cp.Source; source != "/" {
			contents = contents.Directory(source)
		}

		dest := cp.Dest

		if dest != "" && dest[0] == '.' {
			destWorkdir := "/"
			if w := cp.Input.Config.WorkingDir; w != "" {
				destWorkdir = w
			}
			dest = filepath.Join(destWorkdir, dest)
		}

		ct := c.
			Directory(dagger.DirectoryOpts{
				ID: cp.Input.Rootfs.DirectoryID(),
			}).
			WithDirectory(dest, contents, dagger.DirectoryWithDirectoryOpts{
				Include: cp.Include,
				Exclude: cp.Exclude,
			})

		cp.Output.Config = cp.Input.Config
		cp.Output.Platform = cp.Input.Platform
		return cp.Output.Rootfs.SetDirectoryIDBy(ctx, ct)
	})
}
