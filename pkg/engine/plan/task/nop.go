package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Nop{})
}

type Nop struct {
	core.Task
	Input  any `json:"input"`
	Output any `json:"-" wagon:"generated,name=output"`
}

func (e *Nop) Do(ctx context.Context) error {
	e.Output = e.Input
	return nil
}
