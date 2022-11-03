package plan

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	cueerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/tools/flow"
	"dagger.io/dagger"
	"github.com/go-courier/logr"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func NewRunner(root Workdir, value *Value) *Runner {
	r := &Runner{
		root:  root,
		value: value,
	}
	return r
}

type Runner struct {
	root    Workdir
	client  *dagger.Client
	value   *Value
	tasks   []*flow.Task
	target  cue.Path
	targets map[string][]string
}

func (r *Runner) Workdir() Workdir {
	return r.root
}

func printSelectors(w io.Writer, selectors ...cue.Selector) {
	for i, s := range selectors {
		if i > 0 {
			_, _ = fmt.Fprintf(w, ` `)
		}
		_, _ = fmt.Fprintf(w, `%s`, s)
	}
	_, _ = fmt.Fprintf(w, `
`)
}

func (r *Runner) printAllowedTasksTo(w io.Writer) {
	scope := r.target

	_, _ = fmt.Fprintf(w, `
Undefined action:

`)
	printSelectors(w, scope.Selectors()[1:]...)

	_, _ = fmt.Fprintf(w, `
Allowed action:

`)

	taskSelectors := map[string][]cue.Selector{}

	for _, t := range r.tasks {
		selectors := t.Path().Selectors()

		if selectors[0].String() == "actions" {
			publicSelectors := make([]cue.Selector, 0, len(selectors)-1)

			func() {
				for _, selector := range selectors[1:] {
					if selector.String()[0] == '_' {
						return
					}
					publicSelectors = append(publicSelectors, selector)
				}
			}()

			taskSelectors[cue.MakePath(publicSelectors...).String()] = publicSelectors
		}
	}

	keys := make([]string, 0)
	for k := range taskSelectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		printSelectors(w, taskSelectors[k]...)
	}
}

func (r *Runner) resolveDependencies(t *flow.Task) {
	p := t.Path().String()
	if _, ok := r.targets[p]; ok {
		return
	}
	r.targets[p] = make([]string, 0)

	deps := make([]string, 0)
	for _, d := range t.Dependencies() {
		deps = append(deps, d.Path().String())
		r.resolveDependencies(d)
	}

	r.targets[p] = deps
}

func (r *Runner) validateTarget() error {
	r.targets = map[string][]string{}

	for _, t := range r.tasks {
		if strings.HasPrefix(t.Path().String(), r.target.String()) {
			r.resolveDependencies(t)
		}
	}

	if len(r.targets) > 0 {
		//for name, t := range r.targets {
		//	for _, d := range t {
		//		fmt.Println(d, " -> ", name)
		//	}
		//}
		return nil
	}

	buf := bytes.NewBuffer(nil)
	r.printAllowedTasksTo(buf)
	return errors.New(buf.String())
}

func (r *Runner) runTaskFunc(ctx context.Context) flow.TaskFunc {
	taskRunnerFactory := TaskRunnerFactoryContext.From(ctx)

	return func(cueValue cue.Value) (flow.Runner, error) {
		pathKey := cueValue.Path().String()

		if _, ok := r.targets[pathKey]; !ok {
			return nil, nil
		}

		return flow.RunnerFunc(func(t *flow.Task) error {
			tk := WrapTask(t)

			tr, err := taskRunnerFactory.ResolveTaskRunner(tk)
			if err != nil {
				return cueerrors.Wrapf(err, tk.Pos(), "resolve task failed")
			}

			newCtx := logr.WithLogger(ctx,
				logr.FromContext(ctx).
					WithValues("name", pathKey),
			)

			newCtx = WorkdirContext.Inject(newCtx, r.root)
			if err := tr.Run(newCtx); err != nil {
				return cueerrors.Wrapf(err, tk.Pos(), "run task failed")
			}
			return nil
		}), nil
	}
}

func (r *Runner) Run(ctx context.Context, action []string) error {
	actions := append([]string{"actions"}, action...)
	for i := range actions {
		actions[i] = strconv.Quote(actions[i])
	}

	r.target = cue.ParsePath(strings.Join(actions, "."))

	f := flow.New(
		&flow.Config{
			FindHiddenTasks: true,
		},
		r.value.Value(),
		noOpRunner,
	)

	r.tasks = f.Tasks()

	if err := r.validateTarget(); err != nil {
		return err
	}

	return flow.New(
		&flow.Config{
			FindHiddenTasks: true,
		},
		r.value.Value(),
		r.runTaskFunc(ctx),
	).Run(ctx)
}

func noOpRunner(cueValue cue.Value) (flow.Runner, error) {
	v := cueValue.LookupPath(TaskPath)
	if !v.Exists() {
		return nil, nil
	}
	return flow.RunnerFunc(func(t *flow.Task) error {
		return nil
	}), nil
}
