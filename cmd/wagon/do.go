package main

import (
	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/wagon/pkg/engine"
	"github.com/octohelm/wagon/pkg/logutil"
)

func init() {
	cli.AddTo(App, &Do{})
}

type Do struct {
	cli.C
	logutil.Logger
	engine.Engine
	engine.Pipeline
}
