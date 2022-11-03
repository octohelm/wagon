package engine

import (
	"os"

	"github.com/octohelm/wagon/pkg/engine/spec"
	"golang.org/x/net/context"
)

type Build struct {
	Workspaces []string `arg:""`
	// Set target build platform
	Platforms []string `flag:"platform"`
}

func (c *Build) Run(ctx context.Context) error {
	cwd, _ := os.Getwd()

	for i := range c.Workspaces {
		workspace := c.Workspaces[i]

		cc, err := New(WithProjectRoot(workspace), WithWorkdir(cwd))
		if err != nil {
			return err
		}

		for i := range c.Platforms {
			target, err := spec.ParsePlatform(c.Platforms[i])
			if err != nil {
				return err
			}
			if err := cc.Run(ctx, *target); err != nil {
				return err
			}
		}
	}
	return nil
}
