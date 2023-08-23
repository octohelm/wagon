package task

import (
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/logutil"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Print{})
}

type Print struct {
	core.Task
	Input any `json:"input"`
}

func (e *Print) Do(ctx context.Context) error {
	logr.FromContext(ctx).WithValues("input", logutil.CueValue(e.Input)).Info("-")
	return nil
}
