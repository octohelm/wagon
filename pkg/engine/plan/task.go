package plan

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
	"golang.org/x/net/context"
)

func WrapTask(t *flow.Task) Task {
	return &task{task: t}
}

type Task interface {
	Name() string
	Path() cue.Path
	Value() *Value
	Fill(x map[string]any) error
	Iter(ctx context.Context) <-chan Task
}

type task struct {
	task *flow.Task
}

func (t *task) Fill(x map[string]any) error {
	return t.task.Fill(x)
}

func (t *task) Path() cue.Path {
	return t.task.Path()
}

func (t *task) Value() *Value {
	return WrapValue(t.task.Value())
}

func (t *task) Name() string {
	var name string
	if err := t.task.Value().LookupPath(PathForTaskName).Decode(&name); err != nil {
		panic(err)
	}
	return name
}

func (t *task) Iter(ctx context.Context) <-chan Task {
	taskCh := make(chan Task)
	go func() {
		defer close(taskCh)
		select {
		case <-ctx.Done():
			return
		default:
			deps := t.task.Dependencies()
			for i := range deps {
				taskCh <- WrapTask(deps[i])
			}
			taskCh <- t
		}
	}()
	return taskCh
}
