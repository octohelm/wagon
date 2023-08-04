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
	core.DefaultFactory.Register(&ClientFilesystemRead{})
}

type ClientFilesystemRead struct {
	core.Task

	Path     string              `json:"path"`
	Contents core.SecretOrString `json:"contents" wagon:"generated,name=contents"`
}

func (f *ClientFilesystemRead) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		p := f.Path
		if strings.HasPrefix(p, "~") {
			p = fmt.Sprintf("%s%s", os.Getenv("HOME"), p[1:])
		}

		contents, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if secret := f.Contents.Secret; secret != nil {
			fsSecret := c.SetSecret(p, string(contents))

			if err := secret.SetSecretIDBy(ctx, fsSecret); err != nil {
				return err
			}
			f.Contents = core.SecretOrString{
				Secret: secret,
			}
		} else {
			f.Contents = core.SecretOrString{
				Value: string(contents),
			}
		}
		return nil
	})
}
