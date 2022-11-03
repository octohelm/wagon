package fsutil

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type Fs = afero.Fs
type File = afero.File

func NewOsFs() Fs {
	return afero.NewOsFs()
}

func NewBasePathFs(source Fs, path string) Fs {
	if path == "" || path == "." {
		return source
	}
	return &BasePathFs{Source: source, Base: path, Fs: afero.NewBasePathFs(source, path)}
}

type BasePathFs struct {
	Source Fs
	Base   string
	Fs
}

func RealPath(fs Fs) (string, error) {
	base := "."

	for {
		if basePathFs, ok := fs.(*BasePathFs); ok {
			fs = basePathFs.Source
			base = filepath.Join(basePathFs.Base, base)
			continue
		}

		f, err := fs.Stat(".")
		if err != nil {
			return "", err
		}
		base = filepath.Join(f.Name(), base)
		break
	}

	return filepath.Join("/", base), nil
}
