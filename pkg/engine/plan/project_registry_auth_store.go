package plan

import (
	"sync"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/ctxutil"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"golang.org/x/net/context"
)

var RegistryAuthStoreContext = ctxutil.New[RegistryAuthStore]()

type RegistryAuthStore interface {
	Store(address string, auth *Auth)
	ApplyTo(ctx context.Context, container *dagger.Container) *dagger.Container
}

type Auth struct {
	Address  string
	Username string
	SecretID dagger.SecretID
}

func NewRegistryAuthStore() RegistryAuthStore {
	return &registryAuthStore{}
}

type registryAuthStore struct {
	m sync.Map
}

func (r *registryAuthStore) Store(address string, auth *Auth) {
	r.m.Store(address, Auth{
		Address:  address,
		Username: auth.Username,
		SecretID: auth.SecretID,
	})
}

func (r *registryAuthStore) ApplyTo(ctx context.Context, container *dagger.Container) *dagger.Container {
	c := daggerutil.ClientContext.From(ctx)

	r.m.Range(func(key, value any) bool {
		auth := value.(Auth)
		container = container.WithRegistryAuth(auth.Address, auth.Username, c.Secret(auth.SecretID))
		return true
	})

	return container
}
