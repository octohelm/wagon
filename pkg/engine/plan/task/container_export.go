package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
)

func init() {
	core.DefaultFactory.Register(&Export{})
}

type Export struct {
	core.Task

	Input core.FS `json:"input"`
	Dest  string  `json:"dest"`

	ImageConfig *core.ImageConfig `json:"config,omitempty"`
	Platform    string            `json:"platform,omitempty"`
}

func (e *Export) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		fs := plan.WorkdirFor(ctx, plan.WorkdirProject, e.Dest)

		p, err := fsutil.RealPath(fs)
		if err != nil {
			return err
		}

		d := c.Directory(dagger.DirectoryOpts{
			ID: e.Input.DirectoryID(),
		})

		if imageConfig := e.ImageConfig; imageConfig != nil {
			ct := c.Container()

			if v := imageConfig.User; v != "" {
				ct = ct.WithUser(v)
			}

			if v := imageConfig.WorkingDir; v != "" {
				ct = ct.WithWorkdir(v)
			}

			if v := imageConfig.Entrypoint; len(v) > 0 {
				ct = ct.WithEntrypoint(v)
			}

			if v := imageConfig.Cmd; len(v) > 0 {
				ct = ct.WithDefaultArgs(dagger.ContainerWithDefaultArgsOpts{
					Args: v,
				})
			}

			for k, v := range imageConfig.Env {
				ct = ct.WithEnvVariable(k, v)
			}

			for k, v := range imageConfig.Labels {
				ct = ct.WithLabel(k, v)
			}

			_, err = ct.WithDirectory("/", d).Export(ctx, p)
			return err
		}

		_, err = d.Export(ctx, p)
		return err
	})
}
