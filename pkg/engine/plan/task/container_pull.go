package task

import (
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Pull{})
}

type Pull struct {
	core.Task

	Source   string     `json:"source"`
	Platform string     `json:"platform,omitempty" wagon:"generated,name=platform"`
	Auth     *core.Auth `json:"auth,omitempty"`

	Config core.ImageConfig `json:"-" wagon:"generated,name=config"`
	Output core.FS          `json:"-" wagon:"generated,name=output"`

	ResolveMode string `json:"resolveMode,omitempty" wagon:"deprecated"`
}

func (input *Pull) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		ct := c.
			Container(dagger.ContainerOpts{
				Platform: core.DefaultPlatform(input.Platform),
			}).
			From(input.Source)

		ct = plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, ct)
		if a := input.Auth; a != nil {
			ct = a.ApplyTo(ctx, ct, input.Source)
		}

		id, err := ct.ID(ctx)
		if err != nil {
			return err
		}

		if err := input.Config.Resolve(ctx, c, id); err != nil {
			return err
		}
		input.Platform = input.Config.Platform
		return input.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
