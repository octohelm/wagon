package fsutil

import (
	"io/fs"
	"path/filepath"

	"github.com/gobwas/glob"
)

type FileLister interface {
	Range(base string, fn func(filename string) error) error
}

func NewFileLister(f fs.ReadDirFS) FileLister {
	return &fileLister{fs: f}
}

type fileLister struct {
	fs fs.ReadDirFS
}

func (f *fileLister) Range(root string, fn func(filename string) error) error {
	return fs.WalkDir(f.fs, root, func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel := file
		if root != "" && root != "." {
			rel, _ = filepath.Rel(root, file)
		}
		return fn(rel)
	})
}

func WithFilter(lister FileLister, include []string, exclude []string) FileLister {
	globInclude := make([]glob.Glob, len(include))
	globExclude := make([]glob.Glob, len(exclude))

	for i := range include {
		globInclude[i] = glob.MustCompile(include[i])
	}

	for i := range exclude {
		globExclude[i] = glob.MustCompile(exclude[i])
	}

	return &filterLister{
		lister:  lister,
		include: globInclude,
		exclude: globExclude,
	}
}

type filterLister struct {
	include []glob.Glob
	exclude []glob.Glob

	lister FileLister
}

func (f *filterLister) Range(base string, fn func(filename string) error) error {
	return f.lister.Range(base, func(filename string) error {
		if f.Match(filename) {
			return fn(filename)
		}
		return nil
	})
}

func (f *filterLister) Match(filename string) bool {
	for i := range f.exclude {
		if f.exclude[i].Match(filename) {
			return false
		}
	}
	if len(f.include) == 0 {
		return true
	}
	for i := range f.include {
		if f.include[i].Match(filename) {
			return true
		}
	}
	return false
}
