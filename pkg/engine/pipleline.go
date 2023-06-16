package engine

import (
	"bytes"
	cueerrors "cuelang.org/go/cue/errors"
	"fmt"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

type Pipeline struct {
	Action          []string `arg:""`
	ImagePullPrefix string   `flag:",omitempty"`
	Plan            string   `flag:",omitempty" alias:"p"`
	Output          string   `flag:",omitempty" alias:"o"`
}

func (c *Pipeline) SetDefaults() {
	if c.Plan == "" {
		c.Plan = "./wagon.cue"
	}
}

func (c *Pipeline) Run(ctx context.Context) error {
	p, err := New(
		ctx,
		WithPlan(c.Plan),
		WithOutput(c.Output),
	)
	if err != nil {
		return err
	}

	ctx = core.ContextWithImagePullPrefixier(ctx, &imagePullPrefixier{prefix: c.ImagePullPrefix})

	if err := p.Run(ctx, c.Action...); err != nil {
		// print full cue errors if exists
		if errlist := cueerrors.Errors(err); len(errlist) > 0 {
			buf := bytes.NewBuffer(nil)
			for i := range errlist {
				cueerrors.Print(buf, errlist[i], nil)
			}
			return errors.New(buf.String())
		}
		return err
	}
	return nil
}

type imagePullPrefixier struct {
	prefix string
}

func (c *imagePullPrefixier) ImagePullPrefix(name string) string {
	if c.prefix != "" {
		if !strings.HasPrefix(name, c.prefix) {
			if n := len(strings.Split(name, "/")); n <= 2 {
				switch n {
				case 1:
					name = fmt.Sprintf("docker.io/library/%s", name)
				case 2:
					name = fmt.Sprintf("docker.io/%s", name)
				}
			}

			return c.prefix + name
		}
	}
	return name
}
