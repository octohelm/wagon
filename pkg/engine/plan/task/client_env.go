package task

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"dagger.io/dagger"

	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&ClientEnv{})
}

type ClientEnv struct {
	core.Task
	Env map[string]core.SecretOrString `json:",inline" wagon:"generated"`
}

func (s *ClientEnv) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &s.Env)
}

func (s ClientEnv) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Env)
}

func (v *ClientEnv) Do(ctx context.Context) error {
	clientEnvs := getClientEnvs()

	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		for key := range v.Env {
			if envVar, ok := clientEnvs[key]; ok {
				e := v.Env[key]
				if secret := e.Secret; secret != nil {
					envSecret := c.Host().EnvVariable(key).Secret()

					if err := secret.SetSecretIDBy(ctx, envSecret); err != nil {
						return err
					}

					v.Env[key] = core.SecretOrString{
						Secret: secret,
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

func getClientEnvs() map[string]string {
	clientEnvs := map[string]string{}

	for _, i := range os.Environ() {
		parts := strings.SplitN(i, "=", 2)
		if len(parts) == 2 {
			clientEnvs[parts[0]] = parts[1]
		} else {
			clientEnvs[parts[0]] = ""
		}
	}

	return clientEnvs
}
