package task

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Dockerfile{})
}

type Dockerfile struct {
	core.Task

	Source     core.FS           `json:"source"`
	Dockerfile DockerfilePath    `json:"dockerfile"`
	Target     string            `json:"target,omitempty"`
	BuildArg   map[string]string `json:"buildArg"`
	Hosts      map[string]string `json:"hosts"`

	Platform string           `json:"platform,omitempty" wagon:"generated,name=platform"`
	Config   core.ImageConfig `json:"-" wagon:"generated,name=config"`
	Output   core.FS          `json:"-" wagon:"generated,name=output"`

	Auth map[string]core.Auth `json:"auth" wagon:"deprecated"`
}

func (e *Dockerfile) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := c.Directory(dagger.DirectoryOpts{
			ID: e.Source.DirectoryID(),
		})

		buildOpts := dagger.ContainerBuildOpts{
			Dockerfile: e.Dockerfile.Path,
			Target:     e.Target,
		}

		for buildArg := range e.BuildArg {
			buildOpts.BuildArgs = append(buildOpts.BuildArgs, dagger.BuildArg{
				Name:  buildArg,
				Value: e.BuildArg[buildArg],
			})
		}

		ct := c.
			Container(dagger.ContainerOpts{
				Platform: core.DefaultPlatform(e.Platform),
			}).
			Build(dir, buildOpts)

		id, err := ct.ID(ctx)
		if err != nil {
			return err
		}

		if err := e.Config.Resolve(ctx, c, id); err != nil {
			return err
		}

		return e.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}

type DockerfilePath struct {
	Path string `json:"path" default:"Dockerfile"`
}
