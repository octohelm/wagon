package plan

import (
	"github.com/octohelm/wagon/pkg/ctxutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/octohelm/wagon/pkg/engine/spec"
)

var WorkdirContext = ctxutil.New[Workdir]()

type Workdir interface {
	Source(elem ...string) string
	Output(elem ...string) string
	Cache(elem ...string) string
	Pwd() string

	MaskOutput(s string) string
}

func WorkdirFor(pkg spec.Pkg, target spec.Platform, projectRoot string, cwd string) Workdir {
	return &workdir{pkg: pkg, target: target, root: projectRoot, pwd: cwd}
}

type workdir struct {
	pkg    spec.Pkg
	target spec.Platform
	root   string
	pwd    string
}

func (w *workdir) MaskOutput(path string) string {
	return strings.ReplaceAll(path, w.Output()+"/", "${OUTPUT}/")
}

func (w *workdir) Pwd() string {
	return w.pwd
}

func (w *workdir) Cache(elem ...string) string {
	if len(elem) > 0 && filepath.IsAbs(elem[0]) {
		return path.Join(elem...)
	}

	return path.Join(append([]string{
		w.pwd, ".wagon", "cache",
	}, elem...)...)
}

func (w *workdir) Output(elem ...string) string {
	if len(elem) > 0 && filepath.IsAbs(elem[0]) {
		return path.Join(elem...)
	}

	return path.Join(append([]string{
		w.pwd, ".wagon", w.pkg.String(), w.target.StorageKey(),
	}, elem...)...)
}

func (w *workdir) Source(elem ...string) string {
	if len(elem) > 0 && filepath.IsAbs(elem[0]) {
		return path.Join(elem...)
	}

	if filepath.IsAbs(w.root) {
		return path.Join(append([]string{
			w.root,
		}, elem...)...)
	}

	return path.Join(append([]string{
		w.pwd, w.root,
	}, elem...)...)
}
