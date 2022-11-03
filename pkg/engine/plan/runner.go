package plan

import (
	"fmt"
	"github.com/fatih/color"
	"strings"

	"cuelang.org/go/cue"
	cueerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/tools/flow"

	"path/filepath"
	"sync"

	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/spec"
	"golang.org/x/net/context"
)

func NewRunner(pkg spec.Pkg, root Workdir, target spec.Platform, value *Value) *Runner {
	r := &Runner{
		pkg:    pkg,
		root:   root,
		target: target,
		value:  value,
	}
	return r
}

type Runner struct {
	pkg    spec.Pkg
	root   Workdir
	target spec.Platform
	value  *Value

	tasks     sync.Map
	taskQueue []*Task
}

func (r *Runner) Pkg() spec.Pkg {
	return r.pkg
}

func (r *Runner) Workdir() Workdir {
	return r.root
}

func (r *Runner) Target() spec.Platform {
	return r.target
}

func (r *Runner) runTaskFunc(ctx context.Context) flow.TaskFunc {
	wd, _ := filepath.Rel(r.root.Pwd(), r.root.Source())

	resolver := TaskRunnerResolverContext.From(ctx)

	return func(cueValue cue.Value) (flow.Runner, error) {
		v := cueValue.LookupPath(PathForTaskName)
		if !v.Exists() {
			return nil, nil
		}

		return flow.RunnerFunc(func(t *flow.Task) error {
			tt := WrapTask(t)

			for tk := range tt.Iter(ctx) {
				path := tk.Path().String()

				_, ok := r.tasks.Load(path)
				if ok {
					continue
				}
				r.tasks.Store(path, true)

				tr, err := resolver.ResolveTaskRunner(tk)
				if err != nil {
					return cueerrors.Wrapf(err, tk.Value().Pos(), "resolve tk failed")
				}

				newCtx := logr.WithLogger(ctx, logr.FromContext(ctx).WithValues(
					"target", r.target.String(),
					"name", fmt.Sprintf("%s | %s %s", wd, path, color.WhiteString("#%s", strings.ToUpper(tk.Name()))),
				))

				newCtx = WorkdirContext.Inject(newCtx, r.root)

				if err := tr.Run(newCtx); err != nil {
					return err
				}
			}
			return nil
		}), nil
	}
}

func (r *Runner) Run(ctx context.Context) error {
	return flow.New(nil, r.value.Value(), r.runTaskFunc(ctx)).Run(ctx)
}
