package task

import (
	"fmt"
	"os"
	"runtime"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Pull{})
}

type Pull struct {
	core.Task

	Source   string           `json:"source"`
	Platform string           `json:"platform,omitempty" wagon:"generated,name=platform"`
	Config   core.ImageConfig `json:"-" wagon:"generated,name=config"`
	Output   core.FS          `json:"-" wagon:"generated,name=output"`

	ResolveMode string     `json:"resolveMode,omitempty" wagon:"deprecated"`
	Auth        *core.Auth `json:"auth,omitempty" wagon:"deprecated"`
}

func (p *Pull) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		if p.Platform == "" {
			arch := os.Getenv("BUILDKIT_ARCH")
			if arch == "" {
				arch = runtime.GOARCH
			}
			p.Platform = fmt.Sprintf("linux/%s", arch)
		}

		ct := c.
			Container(dagger.ContainerOpts{
				Platform: dagger.Platform(p.Platform),
			}).
			From(p.Source)

		id, err := ct.ID(ctx)
		if err != nil {
			return err
		}

		if err := p.Config.Resolve(ctx, c, id); err != nil {
			return err
		}

		return p.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
