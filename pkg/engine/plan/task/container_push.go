package task

import (
	"context"
	"dagger.io/dagger"
	"encoding/json"

	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/pkg/errors"
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

	Input    core.FS          `json:"input"`
	Config   core.ImageConfig `json:"config"`
	Platform string           `json:"platform,omitempty"`
	Auth     *core.Auth       `json:"auth,omitempty"`

	Result string `json:"-" wagon:"generated,name=result"`
}

func (input *PushImage) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		id := input.Input.DirectoryID()
		if id == "" {
			return errors.Errorf("missing fs")
		}

		dir := c.Directory(dagger.DirectoryOpts{
			ID: id,
		})

		ctr := input.Config.ApplyTo(c.Container(dagger.ContainerOpts{
			Platform: dagger.Platform(input.Platform),
		}).WithRootfs(dir))

		ct := c.Container()

		ct = plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, ct)
		if a := input.Auth; a != nil {
			ct = a.ApplyTo(ctx, ct, input.Dest)
		}

		// prepare before push
		_, err := ctr.Rootfs().Entries(ctx)
		if err != nil {
			return err
		}

		ret, err := ct.Publish(ctx, input.Dest, dagger.ContainerPublishOpts{
			PlatformVariants: []*dagger.Container{
				ctr,
			},
		})
		if err != nil {
			return err
		}
		input.Result = ret
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
	Auth *core.Auth `json:"auth,omitempty"`

	Result string `json:"-" wagon:"generated,name=result"`
}

func (input *PushManifests) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		cts := make([]*dagger.Container, 0, len(input.Inputs))

		for platform := range input.Inputs {
			img := input.Inputs[platform]
			id := img.Input.DirectoryID()
			if id == "" {
				return errors.Errorf("missing fs for %s", platform)
			}

			dir := c.Directory(dagger.DirectoryOpts{
				ID: id,
			})

			ctr := img.Config.ApplyTo(c.Container(dagger.ContainerOpts{Platform: dagger.Platform(platform)}).WithRootfs(dir))

			// prepare before push
			_, err := ctr.Rootfs().Entries(ctx)
			if err != nil {
				return err
			}

			cts = append(cts, ctr)
		}

		ct := plan.RegistryAuthStoreContext.From(ctx).ApplyTo(ctx, c.Container())
		if a := input.Auth; a != nil {
			ct = a.ApplyTo(ctx, ct, input.Dest)
		}

		ret, err := ct.Publish(ctx, input.Dest, dagger.ContainerPublishOpts{
			PlatformVariants: cts,
		})
		if err != nil {
			return err
		}
		input.Result = ret
		return nil
	})
}
