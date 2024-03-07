package core

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"golang.org/x/net/context"
)

func init() {
	DefaultFactory.Register(&Image{})
}

type Image struct {
	Rootfs   FS          `json:"rootfs"`
	Config   ImageConfig `json:"config"`
	Platform string      `json:"platform,omitempty"`
}

func (img *Image) Type() string {
	return "oci"
}

func (img *Image) CanExport() bool {
	return img.Rootfs.CanExport()
}

func (img *Image) ExportTo(ctx context.Context, localPath string) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		rootfs := img.Rootfs.LoadDirectory(c)

		ct := c.Container(dagger.ContainerOpts{
			Platform: DefaultPlatform(img.Platform),
		}).WithRootfs(rootfs)

		ct = img.Config.ApplyTo(ct)

		_, err := ct.Export(ctx, localPath)
		return err
	})
}

type ImageConfig struct {
	WorkingDir string            `json:"workdir" default:""`
	Env        map[string]string `json:"env"`
	Labels     map[string]string `json:"label"`
	Entrypoint []string          `json:"entrypoint"`
	Cmd        []string          `json:"cmd"`
	User       string            `json:"user" default:""`
}

func (p ImageConfig) ApplyTo(c *dagger.Container) *dagger.Container {
	for k := range p.Env {
		c = c.WithEnvVariable(k, p.Env[k])
	}

	for k := range p.Labels {
		c = c.WithLabel(k, p.Labels[k])
	}

	if vv := p.User; vv != "" {
		c = c.WithUser(vv)
	}

	if vv := p.WorkingDir; vv != "" {
		c = c.WithWorkdir(vv)
	}

	if vv := p.Entrypoint; len(vv) != 0 {
		c = c.WithEntrypoint(vv)
	}

	if vv := p.Cmd; len(vv) != 0 {
		c = c.WithDefaultArgs(vv)
	}

	return c
}

func (p ImageConfig) Merge(config ImageConfig) ImageConfig {
	final := ImageConfig{}

	final.WorkingDir = p.WorkingDir
	if config.WorkingDir != "" {
		final.WorkingDir = config.WorkingDir
	}

	final.User = p.User
	if config.WorkingDir != "" {
		final.User = config.User
	}

	final.Entrypoint = p.Entrypoint
	if len(config.Entrypoint) != 0 {
		final.Entrypoint = config.Entrypoint
	}

	final.Cmd = p.Cmd
	if len(config.Cmd) != 0 {
		final.Cmd = config.Cmd
	}

	final.Env = mergeMap(p.Env, config.Env)
	final.Labels = mergeMap(p.Labels, config.Labels)

	return final
}

func (p *ImageConfig) Resolve(ctx context.Context, c *dagger.Client, id dagger.ContainerID) error {
	ret := struct {
		Container daggerutil.Container
	}{}

	err := daggerutil.Query(ctx, c, &ret, fmt.Sprintf(`
query { 
    container(id: %q) {
		id
		entrypoint
		defaultArgs
		workdir
		user
		labels {
			value
			name
		}
		envVariables {
			value
			name
		}
		rootfs {
			id
		}
    }
}
`, id))
	if err != nil {
		return err
	}

	p.WorkingDir = ret.Container.Workdir
	p.Entrypoint = ret.Container.Entrypoint
	p.Cmd = ret.Container.DefaultArgs
	p.User = ret.Container.User

	p.Env = map[string]string{}
	for _, e := range ret.Container.EnvVariables {
		p.Env[e.Name] = e.Value
	}

	p.Labels = map[string]string{}
	for _, e := range ret.Container.Labels {
		p.Labels[e.Name] = e.Value
	}

	return nil
}

func mergeMap[K comparable, V any](a map[K]V, b map[K]V) map[K]V {
	m := make(map[K]V)

	for k := range a {
		m[k] = a[k]
	}

	for k := range b {
		m[k] = b[k]
	}

	return m
}
