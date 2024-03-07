package task

import (
	"context"
	"time"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
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

	Output core.FS `json:"-" wagon:"generated,name=output"`
}

func (e *Exec) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		ct := c.Container().WithRootfs(e.Input.LoadDirectory(c))

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
				ct = ct.WithSecretVariable(k, c.LoadSecretFromID(envVar.Secret.SecretID()))
			} else {
				ct = ct.WithEnvVariable(k, envVar.Value)
			}
		}

		if e.Always {
			// disable cache
			ct = ct.WithEnvVariable("__WAGON_EXEC_STARTED_AT", time.Now().String())
		}

		ct = ct.WithExec(e.Args)

		return e.Output.SetDirectoryIDBy(ctx, ct.Rootfs())
	})
}
