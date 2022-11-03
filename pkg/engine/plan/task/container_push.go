package task

import (
	"encoding/json"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Push{})
}

type Push struct {
	Pusher
}

func (Push) OneOf() []any {
	return []any{
		&PushImage{},
		&PushManifests{},
	}
}

func (p *Push) UnmarshalJSON(bytes []byte) error {
	ret := struct {
		Type string `json:"type"`
	}{}

	if err := json.Unmarshal(bytes, &ret); err != nil {
		return err
	}

	if ret.Type == "manifests" {
		m := &PushManifests{}
		if err := json.Unmarshal(bytes, m); err != nil {
			return err
		}
		p.Pusher = m
	}

	m := &PushImage{}
	if err := json.Unmarshal(bytes, m); err != nil {
		return err
	}
	p.Pusher = m
	return nil
}

type Pusher interface {
	Do(ctx context.Context) error
}

type PushImage struct {
	core.Task

	Dest string `json:"dest"`
	Type string `json:"type" enum:"image"`

	Input  core.FS          `json:"input"`
	Config core.ImageConfig `json:"config"`

	Platform string `json:"platform,omitempty"`
	Result   string `json:"-" wagon:"generated,name=result"`

	Auth *core.Auth `json:"auth,omitempty" wagon:"deprecated"`
}

func (i *PushImage) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		dir := c.Directory(dagger.DirectoryOpts{
			ID: i.Input.DirectoryID(),
		})

		ct := i.Config.ApplyTo(c.Container(dagger.ContainerOpts{
			Platform: dagger.Platform(i.Platform),
		}).WithRootfs(dir))

		ret, err := ct.Publish(ctx, i.Dest)
		if err != nil {
			return err
		}
		i.Result = ret
		return nil
	})

}

type PushManifests struct {
	core.Task

	Dest string `json:"dest"`
	Type string `json:"type" enum:"manifests"`

	Inputs map[string]struct {
		Input  core.FS          `json:"input"`
		Config core.ImageConfig `json:"config"`
	} `json:"inputs"`

	Result string     `json:"-" wagon:"generated,name=result"`
	Auth   *core.Auth `json:"auth,omitempty" wagon:"deprecated"`
}

func (m *PushManifests) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		cts := make([]*dagger.Container, 0)

		for platform := range m.Inputs {
			img := m.Inputs[platform]
			dir := c.Directory(dagger.DirectoryOpts{
				ID: img.Input.DirectoryID(),
			})
			cts = append(cts, img.Config.ApplyTo(c.Container(dagger.ContainerOpts{
				Platform: dagger.Platform(platform),
			}).WithRootfs(dir)))
		}

		ret, err := c.Container().Publish(ctx, m.Dest, dagger.ContainerPublishOpts{
			PlatformVariants: cts,
		})
		if err != nil {
			return err
		}
		m.Result = ret
		return nil
	})
}
