package task

import (
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"golang.org/x/net/context"
)

func init() {
	core.DefaultFactory.Register(&Client{})
}

type Client struct {
	core.Task

	Env map[string]core.SecretOrString `json:"env" wagon:"generated,name=env"`
}

func (v *Client) Do(ctx context.Context) error {
	clientEnvs := map[string]string{}

	for _, i := range os.Environ() {
		parts := strings.SplitN(i, "=", 2)
		if len(parts) == 2 {
			clientEnvs[parts[0]] = parts[1]
		} else {
			clientEnvs[parts[0]] = ""
		}
	}

	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		for key := range v.Env {
			if envVar, ok := clientEnvs[key]; ok {
				e := v.Env[key]
				if s := e.Secret; s != nil {
					if err := s.SetSecretIDBy(ctx, c.Host().EnvVariable(key).Secret()); err != nil {
						return err
					}
					v.Env[key] = core.SecretOrString{
						Secret: s,
					}
				} else {
					v.Env[key] = core.SecretOrString{
						Value: envVar,
					}
				}
			}

		}
		return nil
	})
}
