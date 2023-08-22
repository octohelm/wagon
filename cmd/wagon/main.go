package main

import (
	"os"

	"github.com/innoai-tech/infra/pkg/cli"
	"github.com/octohelm/wagon/internal/version"
	"golang.org/x/net/context"
)

var App = cli.NewApp("wagon", version.Version())

func main() {
	if err := cli.Execute(context.Background(), App, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
