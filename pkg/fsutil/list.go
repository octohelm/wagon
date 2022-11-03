package fsutil

import (
	"github.com/gobwas/glob"
	"io/fs"
	"path/filepath"
)

func NewFinder(include []string, exclude []string) (*Finder, error) {
	f := &Finder{
		Include: make([]glob.Glob, len(include)),
		Exclude: make([]glob.Glob, len(exclude)),
	}

	for i := range include {
		g, err := glob.Compile(include[i])
		if err != nil {
			return nil, err
		}
		f.Include[i] = g
	}

	for i := range exclude {
		g, err := glob.Compile(exclude[i])
		if err != nil {
			return nil, err
		}
		f.Exclude[i] = g
	}

	return f, nil
}

type Finder struct {
	Include []glob.Glob
	Exclude []glob.Glob
}

func (f *Finder) Match(filename string) bool {
	for i := range f.Exclude {
		if f.Exclude[i].Match(filename) {
			return false
		}
	}
	for i := range f.Include {
		if f.Include[i].Match(filename) {
			return true
		}
	}
	return false
}

func (f *Finder) Walk(root string, fn func(file string) error) error {
	return filepath.Walk(root, func(filename string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			if f.Match(filename) {
				rel, _ := filepath.Rel(root, filename)
				return fn(rel)
			}
		}
		return nil
	})
}
