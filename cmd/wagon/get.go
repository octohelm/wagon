package main

import (
	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/wagon/pkg/logutil"
	"golang.org/x/net/context"
	"os"
)

func init() {
	cli.AddTo(App, &Get{})
}

type Get struct {
	cli.C
	logutil.Logger
	GetMod
}

type GetMod struct {
	Pkgs []string `arg:""`

	// Update to latest
	Update bool `flag:"u,omitempty"`

	// declare language for generate. support values: go
	Import string `flag:"i,omitempty"`
}

func (m *GetMod) Run(ctx context.Context) error {
	cwd, _ := os.Getwd()

	c := cuemod.ContextFor(cwd)

	for i := range m.Pkgs {
		p := m.Pkgs[i]

		err := c.Get(
			cuemod.WithOpts(ctx,
				cuemod.OptUpgrade(m.Update),
				cuemod.OptImport(m.Import),
			),
			p,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
