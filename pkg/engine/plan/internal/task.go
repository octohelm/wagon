package internal

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	cueformat "cuelang.org/go/cue/format"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/pkg/errors"
)

type FlowTask interface {
	flowTask()
}

type Task struct {
}

func (Task) flowTask() {

}

func taskFrom(r TaskRegister, t any) *task {
	tpe := reflect.TypeOf(t)
	if tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}

	c := newConvert(r)

	pt := &task{
		tpe:          tpe,
		outputFields: map[string]int{},
		cueType:      c.toCueType(tpe, opt{naming: tpe.Name()}),
	}

	if _, ok := t.(FlowTask); ok {
		pt.flowTask = true
	}

	walkFields(tpe, func(info *fieldInfo) {
		for _, attr := range info.attrs {
			if attr == "generated" {
				pt.outputFields[info.name] = info.idx
			}
		}
	})

	return pt
}

type task struct {
	tpe          reflect.Type
	outputFields map[string]int
	flowTask     bool
	cueType      []byte
}

func (t *task) Name() string {
	return t.tpe.Name()
}

func (t *task) New(planTask plan.Task) (plan.TaskRunner, error) {
	r := &taskRunner{
		task:            planTask,
		inputTaskRunner: reflect.New(t.tpe),
		outputFields:    map[string]int{},
	}

	for f, i := range t.outputFields {
		r.outputFields[f] = i
	}

	return r, nil
}

func (t *task) WriteCueDeclTo(w io.Writer) error {
	b := bytes.NewBuffer(nil)
	name := t.Name()

	if t.flowTask {
		_, _ = fmt.Fprintf(b, `#%s: $wagon: task: name: %q
`, name, name)
	}

	_, _ = fmt.Fprintf(b, `#%s: %s
`, name, t.cueType)

	data, err := cueformat.Source(b.Bytes(), cueformat.Simplify())
	if err != nil {
		return errors.Wrapf(err, `format invalid: %s`, b.Bytes())
	}
	_, err = w.Write(data)
	return err
}
