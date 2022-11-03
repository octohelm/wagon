package internal

import (
	"reflect"
	"time"

	"github.com/octohelm/wagon/pkg/logutil"

	"cuelang.org/go/cue"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type taskRunner struct {
	task plan.Task

	inputTaskRunner reflect.Value
	outputFields    map[string]int
}

func (t *taskRunner) Underlying() any {
	return t.inputTaskRunner.Interface()
}

func (t *taskRunner) Path() cue.Path {
	return t.task.Path()
}

func (t *taskRunner) Task() plan.Task {
	return t.task
}

func (t *taskRunner) Run(ctx context.Context) (e error) {
	inputStepRunner := t.inputTaskRunner.Interface().(plan.StepRunner)

	logr.FromContext(ctx).Info("computing")

	if err := t.task.Decode(inputStepRunner); err != nil {
		return err
	}

	logr.FromContext(ctx).
		WithValues("computed", logutil.CueValue(inputStepRunner)).
		Debug("computed values")

	started := time.Now()
	defer logr.FromContext(ctx).WithValues("duration", time.Now().Sub(started)).Info("computed")

	if err := inputStepRunner.Do(ctx); err != nil {
		return errors.Wrap(err, "do failed")
	}

	values := map[string]any{}

	rv := t.inputTaskRunner

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	for name, i := range t.outputFields {
		if name == "" {
			f := rv.Field(i)
			if f.Kind() == reflect.Map {
				for _, k := range f.MapKeys() {
					key := k.String()
					if key == "$wagon" {
						continue
					}
					values[key] = f.MapIndex(k).Interface()
				}
			}
			continue
		}
		values[name] = rv.Field(i).Interface()
	}

	defer func() {
		if e == nil {
			logr.FromContext(ctx).
				WithValues("result", logutil.CueValue(values)).
				Debug("task result")
		}
	}()

	if err := t.task.Fill(values); err != nil {
		return errors.Wrapf(err, "`%s`: fill results failed", t.task.Path())
	}
	return nil
}
