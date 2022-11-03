package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger.io/dagger"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/logutil"
)

func init() {
	core.DefaultFactory.Register(&Exec{})
}

type Exec struct {
	core.Task

	Input   core.FS                        `json:"input"`
	Mounts  map[string]core.Mount          `json:"mounts"`
	Env     map[string]core.SecretOrString `json:"env"`
	Workdir string                         `json:"workdir" default:"/"`
	Args    []string                       `json:"args"`
	User    string                         `json:"user" default:"root:root"`
	Always  bool                           `json:"always,omitempty"`

	Exit   int     `json:"-" wagon:"generated,name=exit"`
	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (e *Exec) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		ct := c.Container().WithRootfs(c.Directory(dagger.DirectoryOpts{
			ID: e.Input.DirectoryID(),
		}))

		if workdir := e.Workdir; workdir != "" {
			ct = ct.WithWorkdir(workdir)
		}

		if user := e.User; user != "" {
			ct = ct.WithUser(user)
		}

		for n := range e.Mounts {
			ct = e.Mounts[n].MountTo(c, ct)
		}

		for k := range e.Env {
			if envVar := e.Env[k]; envVar.Secret != nil {
				ct = ct.WithSecretVariable(k, c.Secret(envVar.Secret.SecretID()))
			} else {
				ct = ct.WithEnvVariable(k, envVar.Value)
			}
		}

		if e.Always {
			// disable cache
			ct = ct.WithEnvVariable("__WAGON_EXEC_STARTED_AT", time.Now().String())
		}

		ct = ct.WithExec(e.Args)

		_, _ = fmt.Fprint(logutil.Forward(logr.FromContext(ctx).Info), strings.Join(e.Args, " "))

		exitCode, err := ct.ExitCode(ctx)
		if err != nil {
			return err
		}
		e.Exit = exitCode

		return e.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
