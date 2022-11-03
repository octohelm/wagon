package fsutil

import (
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/mod/sumdb/dirhash"
)

func HashDir(fs fs.ReadDirFS, dir string, prefix string, hash dirhash.Hash) (string, error) {
	files := make([]string, 0)

	if err := RangeFile(fs, dir, func(filename string) error {
		files = append(files, filepath.ToSlash(filepath.Join(prefix, filename)))
		return nil
	}); err != nil {
		return "", err
	}

	return hash(files, func(name string) (io.ReadCloser, error) {
		return fs.Open(filepath.Join(dir, strings.TrimPrefix(name, prefix)))
	})
}

func RangeFile(f fs.ReadDirFS, root string, each func(filename string) error) error {
	return fs.WalkDir(f, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := path
		if root != "" && root != "." {
			rel, _ = filepath.Rel(root, path)
		}
		return each(rel)
	})

}
