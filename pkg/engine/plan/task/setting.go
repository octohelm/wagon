package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Setting{})
}

type Setting struct {
	core.Task

	Registry map[string]RegistrySetting `json:"registry"`
}

type RegistrySetting struct {
	Auth core.Auth `json:"auth"`
}

func (v *Setting) Setup() bool {
	return true
}

func (v *Setting) Do(ctx context.Context) error {
	as := plan.RegistryAuthStoreContext.From(ctx)

	for host, r := range v.Registry {
		as.Store(host, &plan.Auth{
			Username: r.Auth.Username,
			SecretID: r.Auth.Secret.SecretID(),
		})
	}

	return nil
}
