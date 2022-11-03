package task

import (
	"dagger.io/dagger"
	"fmt"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/octohelm/wagon/pkg/version/gomod"
	"golang.org/x/net/context"
	"time"
)

func init() {
	core.DefaultFactory.Register(&Source{})
}

type Source struct {
	core.Task

	Path    string   `json:"path" default:"."`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (s *Source) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		fs := plan.WorkdirFor(ctx, plan.WorkdirProject, s.Path)

		p, err := fsutil.RealPath(fs)
		if err != nil {
			return err
		}

		cacheHitsFile := fmt.Sprintf(".fake-cache-hits-%d", time.Now().Unix())

		if version, _ := gomod.LocalRevInfo(p); version != nil {
			cacheHitsFile = cacheHitsFile + "-" + version.Version()
		}

		hostDir := c.Host().Directory(p, dagger.HostDirectoryOpts{
			Include: s.Include,
			Exclude: append(
				s.Exclude,
				cacheHitsFile,
			),
		})

		return s.Output.SetDirectoryIDBy(ctx, hostDir)
	})
}
