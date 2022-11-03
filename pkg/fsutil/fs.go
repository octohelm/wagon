package fsutil

import (
	"os"
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

func OpenFileOrCreate(f Fs, filename string) (File, error) {
	d := filepath.Dir(filename)
	if err := f.MkdirAll(d, 0666); err != nil {
		return nil, err
	}
	file, err := f.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return f.Create(filename)
		}
		return nil, err
	}
	return file, nil
}
