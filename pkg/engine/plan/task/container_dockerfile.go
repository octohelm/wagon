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

	Source     core.FS                 `json:"source"`
	Dockerfile DockerfilePathOrContent `json:"dockerfile"`
	Target     string                  `json:"target,omitempty"`
	BuildArg   map[string]string       `json:"buildArg"`
	Label      map[string]string       `json:"label"`
	Auth       map[string]core.Auth    `json:"auth"`

	Platform string           `json:"platform,omitempty" wagon:"generated,name=platform"`
	Config   core.ImageConfig `json:"-" wagon:"generated,name=config"`
	Output   core.FS          `json:"-" wagon:"generated,name=output"`

	Hosts map[string]string `json:"hosts" wagon:"deprecated"`
}

func (input *Dockerfile) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := input.Source.Directory(c)

		dockerfilePath := input.Dockerfile.Path

		if contents := input.Dockerfile.Contents; contents != "" {
			dockerfilePath = "/Dockerfile"
			dir = dir.WithNewFile(dockerfilePath, contents)
		}

		buildOpts := dagger.ContainerBuildOpts{
			Dockerfile: dockerfilePath,
			Target:     input.Target,
		}

		for buildArg := range input.BuildArg {
			buildArgValue := input.BuildArg[buildArg]

			if buildArgValue == "" {
				continue
			}

			buildOpts.BuildArgs = append(buildOpts.BuildArgs, dagger.BuildArg{
				Name:  buildArg,
				Value: buildArgValue,
			})
		}

		ct := c.Container(dagger.ContainerOpts{
			Platform: core.DefaultPlatform(input.Platform),
		})

		ct = plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, ct)
		for address := range input.Auth {
			ct = input.Auth[address].ApplyTo(ctx, ct, address)
		}

		for label := range input.Label {
			ct = ct.WithLabel(label, input.Label[label])
		}

		ct = ct.Build(dir, buildOpts)

		// hack to trigger Dockerfile build logs
		_, err := ct.Rootfs().Entries(ctx)
		if err != nil {
			return err
		}

		id, err := ct.ID(ctx)
		if err != nil {
			return err
		}
		platform, err := ct.Platform(ctx)
		if err != nil {
			return err
		}

		if err := input.Config.Resolve(ctx, c, id); err != nil {
			return err
		}

		input.Platform = string(platform)
		return input.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}

type DockerfilePathOrContent struct {
	Contents string `json:"contents,omitempty"`
	Path     string `json:"path" default:"Dockerfile"`
}
