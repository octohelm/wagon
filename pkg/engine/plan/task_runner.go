package plan

import (
	"cuelang.org/go/cue"
	"github.com/octohelm/wagon/pkg/ctxutil"
	"golang.org/x/net/context"
)

var PathForTaskName = cue.ParsePath("$wagon.task.name")

type TaskRunner interface {
	Path() cue.Path
	Run(ctx context.Context) error
}

type TaskRunnerResolver interface {
	ResolveTaskRunner(task Task) (TaskRunner, error)
}

var TaskRunnerResolverContext = ctxutil.New[TaskRunnerResolver]()
