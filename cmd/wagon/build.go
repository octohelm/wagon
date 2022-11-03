package main

import (
	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/wagon/pkg/engine"
	"github.com/octohelm/wagon/pkg/logutil"
)

func init() {
	cli.AddTo(App, &Build{})
}

// Pkg kubepkg agent
type Build struct {
	cli.C
	logutil.Logger
	engine.Build
}
