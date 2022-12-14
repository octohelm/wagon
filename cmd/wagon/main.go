package main

import (
	"os"

	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/wagon/pkg/engine"
	"github.com/octohelm/wagon/pkg/version"
	"golang.org/x/net/context"
)

var App = cli.NewApp("wagon", version.Version())

func main() {
	if err := cli.Execute(context.Background(), App, os.Args[1:]); err != nil {
		engine.PrintCueErrorIfNeed(err)
		panic(err)
	}
}
