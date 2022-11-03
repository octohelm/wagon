package cuepkg

import (
	"embed"
	"fmt"
	"github.com/octohelm/cuemod/pkg/cuemod/stdlib"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/octohelm/wagon/pkg/version"
	"github.com/octohelm/wagon/pkg/version/semver"
	"github.com/opencontainers/go-digest"
	"github.com/spf13/afero"
	"golang.org/x/mod/sumdb/dirhash"
	"io/fs"
	"strings"
)

//go:embed dagger.io universe.dagger.io
var daggerPortalModules embed.FS

var (
	DaggerModule         = "dagger.io"
	DaggerUniverseModule = "universe.dagger.io"
	WagonModule          = "wagon.octohelm.tech"
)

func RegistryCueStdlibs() error {
	ver := semver.Parse(version.Version())

	wagonModule, err := createWagonModule()
	if err != nil {
		return err
	}

	if err := registerStdlib(wagonModule, ver, WagonModule); err != nil {
		return nil
	}

	if err := registerStdlib(daggerPortalModules, ver, DaggerModule, DaggerUniverseModule); err != nil {
		return nil
	}

	return nil
}

func registerStdlib(fs fs.ReadDirFS, ver *semver.SemVer, modules ...string) error {
	h, err := fsutil.HashDir(fs, ".", "", dirhash.Hash1)
	if err != nil {
		return err
	}
	stdlib.Register(fs, fmt.Sprintf("%s-20200202235959-%s", ver, strings.ToLower(digest.FromString(h).Hex()[0:12])), modules...)
	return nil
}

func createWagonModule() (fs.ReadDirFS, error) {
	f := afero.NewIOFS(afero.NewMemMapFs())
	file, err := fsutil.OpenFileOrCreate(f, fmt.Sprintf("%s/core/core.cue", WagonModule))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := core.DefaultFactory.WriteCueDeclsTo(file); err != nil {
		return nil, err
	}
	return f, nil
}
