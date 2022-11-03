package internal

import (
	"fmt"
	"io"
	"sort"

	"github.com/pkg/errors"

	"github.com/octohelm/wagon/pkg/engine/plan"
)

type TaskRegister interface {
	Register(t any)
}

type TaskFactory interface {
	TaskRegister
	plan.TaskRunnerFactory

	WriteCueDeclsTo(w io.Writer) error
}

func New() TaskFactory {
	return tasks{}
}

type tasks map[string]*task

func (ts tasks) Register(t any) {
	tk := taskFrom(ts, t)
	ts[tk.Name()] = tk
}

func (ts tasks) ResolveTaskRunner(task plan.Task) (plan.TaskRunner, error) {
	if found, ok := ts[task.Name()]; ok {
		return found.New(task)
	}
	return nil, fmt.Errorf("unknown task `%s`", task.Name())
}

var decls []byte

func (ts tasks) WriteCueDeclsTo(w io.Writer) error {
	if _, err := fmt.Fprint(w, `package core

`); err != nil {
		return err
	}

	if decls == nil {
		names := make([]string, 0, len(ts))
		for k := range ts {
			names = append(names, k)
		}
		sort.Strings(names)

		for i := range names {
			if i > 0 {
				if _, err := fmt.Fprint(w, "\n"); err != nil {
					return err
				}
			}
			if err := ts[names[i]].WriteCueDeclTo(w); err != nil {
				return errors.Wrapf(err, "write decl %s failed", names[i])
			}
		}
	}
	return nil
}
