package old

import (
	"fmt"

	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/shellutil"
	"golang.org/x/net/context"
	"mvdan.cc/sh/v3/interp"
)

func init() {
	//Register("Exec", &execTask{})
}

type execTask struct {
	*plan.Task
}

func (execTask) New(task *plan.Task) plan.TaskRunner {
	return &execTask{Task: task}
}

func (execTask) Decl() []byte {
	return []byte(`
source:  string
dest: string
env:    [Name=string]: _
script: string
`)
}

func (t *execTask) Run(ctx context.Context, planCtx plan.Context) error {
	l := logr.FromContext(ctx)
	val := t.Value()
	workdir := planCtx.Workdir()

	var script string
	var envVars shellutil.EnvVars

	if err := val.Lookup("script").Decode(&script); err != nil {
		return err
	}

	if err := val.Lookup("env").Decode(&envVars); err != nil {
		return err
	}

	l.Info(fmt.Sprintf("%s%s", envVars.String(), workdir.MaskOutput(script)))

	if err := shellutil.Exec(
		ctx, envVars.String()+script,
		interp.Dir(workdir.Source()),
	); err != nil {
		return err
	}

	return nil
}
