package task

import (
	"context"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&ImageSet{})
}

type ImageSet struct {
	core.Task
	Input  core.Image       `json:"input"`
	Config core.ImageConfig `json:"config"`
	Output core.Image       `json:"-" wagon:"generated,name=output"`
}

func (e *ImageSet) Do(ctx context.Context) error {
	e.Output.Rootfs = e.Input.Rootfs
	e.Output.Platform = e.Input.Platform
	e.Output.Config = e.Input.Config.Merge(e.Config)
	return nil
}
