package task

import (
	"fmt"
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

	Env        map[string]core.SecretOrString `json:"env" wagon:"generated,name=env"`
	Filesystem map[string]struct {
		Read  *ClientFile `json:"read,omitempty"`
		Write *ClientFile `json:"write,omitempty" wagon:"deprecated"`
	} `json:"filesystem" wagon:"generated,name=filesystem"`
}

type ClientFile struct {
	Path     string              `json:"path,omitempty"`
	Contents core.SecretOrString `json:"contents"`
}

func (v *Client) Do(ctx context.Context) error {
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

		for path := range v.Filesystem {
			fs := v.Filesystem[path]
			if cf := fs.Read; cf != nil {
				p := cf.Path
				if p == "" {
					p = path
				}

				if strings.HasPrefix(p, "~") {
					p = fmt.Sprintf("%s%s", os.Getenv("HOME"), p[1:])
				}

				contents, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				if secret := cf.Contents.Secret; secret != nil {
					fsSecret := c.Directory().WithNewFile(p, string(contents)).File(p).Secret()
					if err := secret.SetSecretIDBy(ctx, fsSecret); err != nil {
						return err
					}
					*v.Filesystem[path].Read = ClientFile{
						Path: path,
						Contents: core.SecretOrString{
							Secret: secret,
						},
					}
				} else {
					*v.Filesystem[path].Read = ClientFile{
						Path: path,
						Contents: core.SecretOrString{
							Value: string(contents),
						},
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
