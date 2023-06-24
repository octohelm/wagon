package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&CloneWithoutNull{})
}

type CloneWithoutNull struct {
	core.Task
	Input  any `json:"input"`
	Output any `json:"-" wagon:"generated,name=output"`
}

func (e *CloneWithoutNull) Do(ctx context.Context) error {
	e.Output = cloneWithoutNull(e.Input)
	return nil
}

func cloneWithoutNull(v any) any {
	switch x := v.(type) {
	case map[string]any:
		vv := map[string]any{}
		for k := range x {
			if pv := cloneWithoutNull(x[k]); pv != nil {
				vv[k] = pv
			}
		}
		return vv
	case []any:
		s := make([]any, 0, len(x))
		for i := range x {
			if ev := cloneWithoutNull(x[i]); ev != nil {
				s = append(s, ev)
			}
		}
		return s
	}
	return v
}
