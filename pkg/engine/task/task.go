package task

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-courier/logr"
	"github.com/opencontainers/go-digest"
	"go/ast"
	"io"
	"reflect"
	"sort"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/octohelm/wagon/pkg/engine/plan"
)

var Default = tasks{}

type Runner interface {
	DigestInputs(x any) (digest.Digest, error)
	Do(ctx context.Context, dgst digest.Digest) (Omit, error)
}

type PreRunner interface {
	PreDo(ctx context.Context) error
}

type Omit = func(t plan.Task) error

func Register(runner Runner) {
	Default.Register(runner)
}

type task struct {
	tpe reflect.Type
}

func (t *task) Name() string {
	return t.tpe.Name()
}

func (t *task) New(x plan.Task) (plan.TaskRunner, error) {
	tr := reflect.New(t.tpe).Interface().(Runner)
	return &taskRunner{x: x, runner: tr}, nil
}

type taskRunner struct {
	x      plan.Task
	runner Runner
}

func (t *taskRunner) Path() cue.Path {
	return t.x.Path()
}

func (t *taskRunner) Run(ctx context.Context) (e error) {
	if err := t.x.Value().Decode(t.runner); err != nil {
		return err
	}

	if preRunner, ok := t.runner.(PreRunner); ok {
		if err := preRunner.PreDo(ctx); err != nil {
			return err
		}
	}

	dgst, err := t.runner.DigestInputs(t.runner)
	if err != nil {
		return err
	}

	l := logr.FromContext(ctx)
	l.Info("started.")
	startedAt := time.Now()
	defer func() {
		if e == nil {
			l.WithValues("cost", time.Now().Sub(startedAt)).Info(
				fmt.Sprintf("done. digest %s cached.", dgst),
			)
		}
	}()

	omit, err := t.runner.Do(ctx, dgst)
	if err != nil {
		return err
	}

	if omit != nil {
		return omit(t.x)
	}

	return nil
}

func toCueType(tpe reflect.Type) string {
	if tpe.PkgPath() != "" {
		return "#" + tpe.Name()
	}

	switch tpe.Kind() {
	case reflect.Slice:
		return fmt.Sprintf("[...%s]", toCueType(tpe.Elem()))
	default:
		return tpe.String()
	}
}

func walkField(w io.Writer, s reflect.Type) {
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !ast.IsExported(f.Name) {
			continue
		}

		ft := toCueType(f.Type)

		if jsonTag, ok := f.Tag.Lookup("json"); ok {
			name := strings.Split(jsonTag, ",")[0]
			if name == "" {
				name = f.Name
			}

			_, _ = fmt.Fprintf(w, `%s: %s`, name, ft)

			if defaultValue, ok := f.Tag.Lookup("default"); ok {
				switch ft {
				case "[]byte":
					_, _ = fmt.Fprintf(w, ` | *'%s'`, defaultValue)
				case "string":
					_, _ = fmt.Fprintf(w, ` | *%q`, defaultValue)
				default:
					_, _ = fmt.Fprintf(w, ` | *%v`, defaultValue)
				}
			}

			_, _ = fmt.Fprintf(w, "\n")
		}
	}
}

func (t *task) Decl() []byte {
	b := bytes.NewBuffer(nil)

	name := t.Name()

	_, _ = fmt.Fprintf(b, `#%s: {`, name)

	_, _ = fmt.Fprintf(b, `
$wagon: task: name: %q
`, name)

	walkField(b, t.tpe)

	_, _ = fmt.Fprintf(b, `}

`)

	return b.Bytes()
}

type tasks map[string]*task

func (ts tasks) Register(t Runner) {
	tpe := reflect.TypeOf(t)
	if tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}
	tk := &task{tpe: tpe}
	ts[tk.Name()] = tk
}

func (ts tasks) ResolveTaskRunner(task plan.Task) (plan.TaskRunner, error) {
	if found, ok := ts[task.Name()]; ok {
		return found.New(task)
	}
	return nil, fmt.Errorf("unknown task `%s`", task.Name())
}

func (ts tasks) AllDefs() []byte {
	b := bytes.NewBuffer(nil)

	names := make([]string, 0, len(ts))
	for k := range ts {
		names = append(names, k)
	}
	sort.Strings(names)

	for i := range names {
		b.Write(ts[names[i]].Decl())
	}

	return b.Bytes()
}
