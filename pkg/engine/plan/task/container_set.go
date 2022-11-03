package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Set{})
}

type Set struct {
	core.Task

	Input  core.ImageConfig `json:"input"`
	Config core.ImageConfig `json:"config"`

	Output core.ImageConfig `json:"-" wagon:"generated,name=output"`
}

func (e *Set) Do(ctx context.Context) error {
	e.Output = e.Input.Merge(e.Config)
	return nil
}
