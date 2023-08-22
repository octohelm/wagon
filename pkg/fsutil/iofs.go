package fsutil

import (
	"io/fs"
	"path/filepath"
)

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
