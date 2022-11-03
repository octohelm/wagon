package task

import (
	"context"
	"github.com/octohelm/wagon/pkg/engine/plan"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Dockerfile{})
}

type Dockerfile struct {
	core.Task

	Source     core.FS              `json:"source"`
	Dockerfile DockerfilePath       `json:"dockerfile"`
	Target     string               `json:"target,omitempty"`
	BuildArg   map[string]string    `json:"buildArg"`
	Label      map[string]string    `json:"label"`
	Auth       map[string]core.Auth `json:"auth"`

	Platform string           `json:"platform,omitempty" wagon:"generated,name=platform"`
	Config   core.ImageConfig `json:"-" wagon:"generated,name=config"`
	Output   core.FS          `json:"-" wagon:"generated,name=output"`

	Hosts map[string]string `json:"hosts" wagon:"deprecated"`
}

func (input *Dockerfile) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := c.Directory(dagger.DirectoryOpts{
			ID: input.Source.DirectoryID(),
		})

		buildOpts := dagger.ContainerBuildOpts{
			Dockerfile: input.Dockerfile.Path,
			Target:     input.Target,
		}

		for buildArg := range input.BuildArg {
			buildOpts.BuildArgs = append(buildOpts.BuildArgs, dagger.BuildArg{
				Name:  buildArg,
				Value: input.BuildArg[buildArg],
			})
		}

		ct := c.Container(dagger.ContainerOpts{
			Platform: core.DefaultPlatform(input.Platform),
		})

		for label := range input.Label {
			ct = ct.WithLabel(label, input.Label[label])
		}

		ct = plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, ct)
		for address := range input.Auth {
			ct = input.Auth[address].ApplyTo(ctx, ct, address)
		}

		ct = ct.Build(dir, buildOpts)

		if _, err := ct.ExitCode(ctx); err != nil {
			return err
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

type DockerfilePath struct {
	Path string `json:"path" default:"Dockerfile"`
}
