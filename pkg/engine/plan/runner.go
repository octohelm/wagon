package plan

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	cueerrors "cuelang.org/go/cue/errors"
	"cuelang.org/go/tools/flow"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/octohelm/wagon/pkg/logutil"
	"github.com/pkg/errors"
)

func NewRunner(value *Value, output string, exporters ...daggerutil.Exporter) *Runner {
	r := &Runner{
		input:     value,
		exporters: exporters,
		output:    output,
	}
	return r
}

type Runner struct {
	input *Value

	target    cue.Path
	output    string
	exporters []daggerutil.Exporter

	setups  map[string][]string
	targets map[string][]string
}

func (r *Runner) printAllowedTasksTo(w io.Writer, tasks []*flow.Task) {
	scope := r.target

	_, _ = fmt.Fprintf(w, `
Undefined action:

`)
	printSelectors(w, scope.Selectors()[1:]...)

	_, _ = fmt.Fprintf(w, `
Allowed action:

`)

	taskSelectors := map[string][]cue.Selector{}

	for _, t := range tasks {
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

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.TabIndent)
	defer func() {
		_ = tw.Flush()
	}()

	for _, k := range keys {
		printSelectors(tw, taskSelectors[k]...)

		taskValue := r.input.Value().LookupPath(cue.ParsePath("actions." + k))

		if v := taskValue.LookupPath(cue.ParsePath("output")); v.Exists() {
			if v.LookupPath(cue.ParsePath("$wagon.fs")).Exists() {
				_, _ = fmt.Fprintf(tw, "\t\t[--output=<dir>]")
			} else if v.LookupPath(cue.ParsePath("rootfs.$wagon.fs")).Exists() {
				_, _ = fmt.Fprintf(tw, "\t\t[--output=<path_to_oci_tar>]")
			}
		}

		if n := taskValue.Source(); n != nil {
			for _, c := range ast.Comments(n) {
				_, _ = fmt.Fprintf(tw, "\t\t%s", strings.TrimSpace(c.Text()))
			}
		}

		_, _ = fmt.Fprintln(tw)
	}
}

func printSelectors(w io.Writer, selectors ...cue.Selector) {
	for i, s := range selectors {
		if i > 0 {
			_, _ = fmt.Fprintf(w, ` `)
		}
		_, _ = fmt.Fprintf(w, `%s`, s.String())
	}
}

func (r *Runner) resolveDependencies(t *flow.Task, collection map[string][]string) {
	p := t.Path().String()
	if _, ok := collection[p]; ok {
		return
	}
	// avoid cycle
	collection[p] = make([]string, 0)

	deps := make([]string, 0)
	for _, d := range t.Dependencies() {
		deps = append(deps, d.Path().String())
		r.resolveDependencies(d, collection)
	}

	collection[p] = deps
}

func (r *Runner) prepareTasks(ctx context.Context, tasks []*flow.Task) error {
	taskRunnerFactory := TaskRunnerFactoryContext.From(ctx)

	r.setups = map[string][]string{}
	r.targets = map[string][]string{}

	for i := range tasks {
		tk := WrapTask(tasks[i])

		t, err := taskRunnerFactory.ResolveTaskRunner(tk)
		if err != nil {
			return cueerrors.Wrapf(err, tk.Pos(), "resolve task failed")
		}

		if _, ok := t.Underlying().(interface{ Setup() bool }); ok {
			r.resolveDependencies(tasks[i], r.setups)
		}

		if strings.HasPrefix(tk.Path().String(), r.target.String()) {
			r.resolveDependencies(tasks[i], r.targets)
		}
	}

	if r.target.String() != "actions" && len(r.targets) > 0 {
		if os.Getenv("WAGON_GRAPH") != "" {
			fmt.Println(printGraph(r.targets))
		}
		return nil
	}

	buf := bytes.NewBuffer(nil)
	r.printAllowedTasksTo(buf, tasks)
	return errors.New(buf.String())
}

func (r *Runner) runTaskFunc(taskRunnerFactory TaskRunnerFactory, shouldRun func(p cue.Path) bool) flow.TaskFunc {
	sessionToken := os.Getenv("DAGGER_SESSION_TOKEN")

	return func(cueValue cue.Value) (flow.Runner, error) {
		p := cueValue.Path()
		if !(shouldRun(cueValue.Path())) {
			return nil, nil
		}

		return flow.RunnerFunc(func(t *flow.Task) error {
			tk := WrapTask(t)

			tr, err := taskRunnerFactory.ResolveTaskRunner(tk)
			if err != nil {
				return cueerrors.Wrapf(err, tk.Pos(), "resolve task failed")
			}

			displayName := fmt.Sprintf("%s #%s", p, tk.Name())

			c := t.Context()

			c = logr.WithLogger(
				c,
				logr.FromContext(c).WithValues("name", displayName),
			)

			c = daggerutil.ClientContext.Inject(
				c,
				daggerutil.ClientContext.From(c).
					Pipeline(
						fmt.Sprintf("%s%s%s",
							daggerutil.PipelinePrefix, sessionToken, displayName,
						),
					),
			)

			if err := tr.Run(c); err != nil {
				return cueerrors.Wrapf(err, tk.Pos(), "exec task failed")
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

	cueValue := r.input.Value()

	r.target = cue.ParsePath(strings.Join(actions, "."))

	f := flow.New(
		&flow.Config{
			FindHiddenTasks: true,
		},
		cueValue,
		noOpRunner,
	)

	if err := r.prepareTasks(ctx, f.Tasks()); err != nil {
		return err
	}

	return daggerutil.ConnectDo(ctx, func(ctx context.Context) error {
		logr.FromContext(ctx).WithValues("name", "Pipeline").Debug("starting...")

		preparedCueValue, err := r.exec(ctx, cueValue, func(p cue.Path) bool {
			_, ok := r.setups[p.String()]
			return ok
		})

		if err != nil {
			return err
		}

		ret, err := r.exec(ctx, preparedCueValue, func(p cue.Path) bool {
			_, ok := r.targets[p.String()]
			return ok
		})

		if err != nil {
			return err
		}

		l := logr.FromContext(ctx)

		if o := ret.LookupPath(r.target).LookupPath(cue.ParsePath("result")); o.Exists() {
			out := o.Value()
			l.WithValues("result", logutil.CueValue(out)).Info("done")

			return nil
		}

		if o := ret.LookupPath(r.target).LookupPath(cue.ParsePath("output")); o.Exists() {
			l.WithValues("output", logutil.CueValue(o.Value())).Info("done")

			if r.output != "" {
				l = l.WithValues("name", "Export", "dest", r.output)

				for i := range r.exporters {
					e := r.exporters[i]

					_ = o.Value().Decode(e)

					if e.CanExport() {
						fs := WorkdirFor(ctx, WorkdirProject, r.output)
						dest, err := fsutil.RealPath(fs)
						if err != nil {
							return err
						}
						return e.ExportTo(logr.WithLogger(ctx, l), dest)
					}
				}
			}
		}

		return nil
	})
}

func (r *Runner) exec(ctx context.Context, cueValue cue.Value, shouldRun func(p cue.Path) bool) (cue.Value, error) {
	f := flow.New(
		&flow.Config{
			FindHiddenTasks: true,
		},
		cueValue,
		r.runTaskFunc(TaskRunnerFactoryContext.From(ctx), shouldRun),
	)

	if err := f.Run(ctx); err != nil {
		return cue.Value{}, err
	}

	return f.Value(), nil
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

func printGraph(targets map[string][]string) (string, error) {
	buffer := bytes.NewBuffer(nil)

	w, err := zlib.NewWriterLevel(buffer, 9)
	if err != nil {
		return "", errors.Wrap(err, "fail to create the w")
	}

	_, _ = fmt.Fprintf(w, "direction: right\n")
	for name, deps := range targets {
		for _, d := range deps {
			_, _ = fmt.Fprintf(w, "%q -> %q\n", d, name)
		}
	}
	_ = w.Close()
	if err != nil {
		return "", errors.Wrap(err, "fail to create the payload")
	}
	return fmt.Sprintf("https://kroki.io/d2/svg/%s?theme=101", base64.URLEncoding.EncodeToString(buffer.Bytes())), nil
}
