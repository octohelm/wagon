package main

import (
	"context"
	"os"

	"github.com/octohelm/wagon/pkg/version"

	"github.com/octohelm/wagon/pkg/logutil"

	"github.com/innoai-tech/infra/devpkg/gengo"
	"github.com/innoai-tech/infra/pkg/cli"
	_ "github.com/octohelm/courier/devpkg/clientgen"
	_ "github.com/octohelm/courier/devpkg/operatorgen"
	_ "github.com/octohelm/gengo/devpkg/deepcopygen"
	_ "github.com/octohelm/gengo/devpkg/runtimedocgen"
	_ "github.com/octohelm/storage/devpkg/enumgen"
	_ "github.com/octohelm/storage/devpkg/tablegen"
)

var App = cli.NewApp("gengo", version.Version())

func init() {
	cli.AddTo(App, &struct {
		cli.C `name:"gen"`
		logutil.Logger
		gengo.Gengo
	}{})
}

func main() {
	if err := cli.Execute(context.Background(), App, os.Args[1:]); err != nil {
		panic(err)
	}
}
