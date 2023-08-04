package core

import (
	"context"
	"sync"

	"dagger.io/dagger"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

func init() {
	DefaultFactory.Register(&FS{})
}

type FS struct {
	Meta struct {
		Fs struct {
			ID string `json:"id,omitempty"`
		} `json:"fs"`
	} `json:"$wagon"`
}

func (fs *FS) SetDirectoryIDBy(ctx context.Context, d *dagger.Directory) error {
	dir, err := d.Sync(ctx)
	if err != nil {
		return err
	}
	id, err := dir.ID(ctx)
	if err != nil || id == "" {
		return errors.Wrap(err, "resolve dir id failed")
	}
	fs.SetDirectoryID(id)
	return nil
}

var fsids = sync.Map{}

func (fs *FS) SetDirectoryID(id dagger.DirectoryID) {
	key := digest.FromString(string(id)).String()
	fsids.Store(key, id)
	fs.Meta.Fs.ID = key
}

func (fs *FS) DirectoryID() dagger.DirectoryID {
	if k, ok := fsids.Load(fs.Meta.Fs.ID); ok {
		return k.(dagger.DirectoryID)
	}
	return ""
}

func (f *FS) Type() string {
	return "fs"
}

func (f *FS) CanExport() bool {
	return f.Meta.Fs.ID != ""
}

func (f *FS) ExportTo(ctx context.Context, localPath string) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		rootfs := c.Directory(dagger.DirectoryOpts{
			ID: f.DirectoryID(),
		})
		_, err := rootfs.Export(ctx, localPath)
		return err
	})
}
