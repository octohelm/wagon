package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/octohelm/wagon/pkg/version/gomod"
	"github.com/pkg/errors"
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
	w := plan.WorkdirFor(ctx, plan.WorkdirProject)
	p, err := fsutil.RealPath(w)
	if err != nil {
		return err
	}

	revInfo, err := gomod.LocalRevInfo(p)
	if err != nil {
		return errors.Wrapf(err, "load version failed")
	}

	v.Output = revInfo.Version()
	return nil
}
