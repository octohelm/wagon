package cuepkg

import (
	"embed"
	"fmt"
	"github.com/octohelm/cuemod/pkg/cuemod/stdlib"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/octohelm/wagon/pkg/version"
	"github.com/octohelm/wagon/pkg/version/semver"
	"github.com/spf13/afero"
	"io"
	"io/fs"
)

//go:embed dagger.io universe.dagger.io wagon.octohelm.tech
var daggerPortalModules embed.FS

var (
	WagonModule = "wagon.octohelm.tech"

	DaggerModule         = "dagger.io"
	DaggerUniverseModule = "universe.dagger.io"
)

func RegistryCueStdlibs() error {
	wagonModule, err := createWagonModule(daggerPortalModules)
	if err != nil {
		return err
	}

	ver := semver.Parse(version.Version())

	if err := registerStdlib(wagonModule, ver, WagonModule, DaggerModule, DaggerUniverseModule); err != nil {
		return nil
	}

	return nil
}

func registerStdlib(fs fs.ReadDirFS, ver *semver.SemVer, modules ...string) error {
	stdlib.Register(fs, ver.String(), modules...)
	return nil
}

func createWagonModule(otherFs ...fs.ReadDirFS) (fs.ReadDirFS, error) {
	mfs := afero.NewMemMapFs()

	for i := range otherFs {
		f := otherFs[i]
		if err := fsutil.RangeFile(f, ".", func(filename string) error {
			file, err := f.Open(filename)
			if err != nil {
				return err
			}
			defer file.Close()
			newFile, err := fsutil.OpenFileOrCreate(mfs, filename)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(newFile, file); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	file, err := fsutil.OpenFileOrCreate(mfs, fmt.Sprintf("%s/core/core.cue", WagonModule))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if err := core.DefaultFactory.WriteCueDeclsTo(file); err != nil {
		return nil, err
	}

	return afero.NewIOFS(mfs), nil
}
