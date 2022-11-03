package plan

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/tools/flow"
	"fmt"
	"golang.org/x/net/context"
	"os"
)

func WrapTask(t *flow.Task) Task {
	name, _ := t.Value().LookupPath(TaskPath).String()
	return &task{task: t, name: name}
}

type Task interface {
	Name() string
	Pos() token.Pos
	Path() cue.Path
	Decode(input any) error
	Fill(values map[string]any) error
}

type StepRunner interface {
	Do(ctx context.Context) error
}

type Embed interface {
	EmbedPayload() any
}

type task struct {
	name string
	task *flow.Task
}

func (t *task) Pos() token.Pos {
	return t.task.Value().Pos()
}

func (t *task) Decode(input any) error {
	v := t.task.Value()

	if err := v.Decode(input); err != nil {
		_, _ = fmt.Fprint(os.Stdout, t.Value().Source())
		_, _ = fmt.Fprintln(os.Stdout)
		return err
	}

	return nil
}

func (t *task) Name() string {
	return t.name
}

func (t *task) Path() cue.Path {
	return t.task.Path()
}

func (t *task) Value() *Value {
	return WrapValue(t.task.Value())
}

func (t *task) Fill(values map[string]any) error {
	return t.task.Fill(values)
}
