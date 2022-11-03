package fsutil

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"golang.org/x/mod/sumdb/dirhash"
)

type IOFS = afero.IOFS

func OpenFileOrCreate(fs IOFS, filename string) (File, error) {
	d := filepath.Dir(filename)
	if err := fs.MkdirAll(d, 0666); err != nil {
		return nil, err
	}
	f, err := fs.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return fs.Create(filename)
		}
		return nil, err
	}
	return f, nil
}

func HashDir(fs fs.ReadDirFS, dir string, prefix string, hash dirhash.Hash) (string, error) {
	ls := NewFileLister(fs)

	files := make([]string, 0)

	if err := ls.Range(dir, func(filename string) error {
		files = append(files, filepath.ToSlash(filepath.Join(prefix, filename)))
		return nil
	}); err != nil {
		return "", err
	}

	return hash(files, func(name string) (io.ReadCloser, error) {
		return fs.Open(filepath.Join(dir, strings.TrimPrefix(name, prefix)))
	})
}
