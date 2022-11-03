package plan

import (
	"path"

	"github.com/octohelm/wagon/pkg/ctxutil"
	"github.com/octohelm/wagon/pkg/fsutil"
	"golang.org/x/net/context"
)

var WorkdirContext = ctxutil.New[Workdir]()

func WorkdirFor(ctx context.Context, typ WorkdirType, elem ...string) fsutil.Fs {
	return WorkdirContext.From(ctx).Fs(typ, elem...)
}

type WorkdirType string

const (
	WorkdirRoot     WorkdirType = "root"
	WorkdirProject  WorkdirType = "project"
	WorkdirSnapshot WorkdirType = "snapshot"
	WorkdirCache    WorkdirType = "cache"
)

type Workdir interface {
	Fs(typ WorkdirType, elem ...string) fsutil.Fs
}

func NewWorkdir(pwd string, base string) Workdir {
	rootfs := fsutil.NewOsFs()

	return &workdir{
		rootfs: rootfs,
		pwdfs:  fsutil.NewBasePathFs(rootfs, pwd),
		base:   base,
	}
}

type workdir struct {
	rootfs fsutil.Fs
	pwdfs  fsutil.Fs
	base   string
}

func (w *workdir) Fs(tpe WorkdirType, elem ...string) fsutil.Fs {
	switch tpe {
	case WorkdirRoot:
		return w.rootfs
	case WorkdirProject:
		return fsutil.NewBasePathFs(w.pwdfs, path.Join(append([]string{w.base}, elem...)...))
	case WorkdirSnapshot:
		return fsutil.NewBasePathFs(w.pwdfs, path.Join(append([]string{".wagon", "snapshots"}, elem...)...))
	case WorkdirCache:
		return fsutil.NewBasePathFs(w.pwdfs, path.Join(append([]string{".wagon", "cache"}, elem...)...))
	}
	return fsutil.NewBasePathFs(w.pwdfs, path.Join(elem...))
}
