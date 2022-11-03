package core

import (
	"context"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
)

type Auth struct {
	Username string `json:"username"`
	Secret   Secret `json:"secret"`
}

func (a Auth) ApplyTo(ctx context.Context, ct *dagger.Container, address string) *dagger.Container {
	c := daggerutil.ClientContext.From(ctx)
	return ct.WithRegistryAuth(address, a.Username, c.Secret(a.Secret.SecretID()))
}
