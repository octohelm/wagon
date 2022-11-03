package task

import (
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&ImagePull{})
}

type ImagePull struct {
	core.Task

	Source   string `json:"source"`
	Platform string `json:"platform,omitempty"`

	Output core.Image `json:"-" wagon:"generated,name=output"`

	ResolveMode string     `json:"resolveMode,omitempty" wagon:"deprecated"`
	Auth        *core.Auth `json:"auth,omitempty" wagon:"deprecated"`
}

func (p *ImagePull) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		ct := c.
			Container(dagger.ContainerOpts{
				Platform: core.DefaultPlatform(p.Platform),
			}).
			From(p.Source)

		id, err := ct.ID(ctx)
		if err != nil {
			return err
		}
		if err := p.Output.Config.Resolve(ctx, c, id); err != nil {
			return err
		}
		p.Output.Platform = p.Platform
		return p.Output.Rootfs.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
