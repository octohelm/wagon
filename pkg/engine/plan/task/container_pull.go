package task

import (
	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/pkg/errors"
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
		ct := c.Container(dagger.ContainerOpts{
			Platform: core.DefaultPlatform(input.Platform),
		})

		ct = plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, ct)
		if a := input.Auth; a != nil {
			ct = a.ApplyTo(ctx, ct, input.Source)
		}

		ct = ct.From(input.Source)

		id, err := ct.ID(ctx)
		if err != nil {
			return errors.Wrapf(err, "Pull %s failed.", input.Source)
		}

		platform, err := ct.Platform(ctx)
		if err != nil {
			return errors.Wrapf(err, "Resolve Platform %s failed.", input.Source)
		}

		if err := input.Config.Resolve(ctx, c, id); err != nil {
			return err
		}

		input.Platform = string(platform)
		return input.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
