package task

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"

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
	if err := json.Unmarshal(data, &s.Env); err != nil {
		return err
	}
	delete(s.Env, "$wagon")
	return nil
}

func (s *ClientEnv) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Env)
}

func (v *ClientEnv) Do(ctx context.Context) error {
	clientEnvs := getClientEnvs()

	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		env := map[string]core.SecretOrString{}

		for key := range v.Env {
			e := v.Env[key]

			if envVar, ok := clientEnvs[key]; ok {
				if secret := e.Secret; secret != nil {
					envSecret := c.Host().EnvVariable(key).Secret()

					if err := secret.SetSecretIDBy(ctx, envSecret); err != nil {
						return err
					}

					env[key] = core.SecretOrString{
						Secret: secret,
					}
				} else {
					env[key] = core.SecretOrString{
						Value: envVar,
					}
				}
			} else {
				if secret := e.Secret; secret != nil {
					return errors.Errorf("EnvVar %s is not defined.", key)
				}
				env[key] = e
			}
		}

		v.Env = env

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
