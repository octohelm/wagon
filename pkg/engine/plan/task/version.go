package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Version{})
}

type Version struct {
	core.Task

	Output string `json:"-" wagon:"generated,name=output"`
}

func (v *Version) Do(ctx context.Context) error {
	v.Output = plan.MetaContext.From(ctx).Version
	return nil
}
