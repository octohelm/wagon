package core

import (
	"context"
	"sync"

	"github.com/opencontainers/go-digest"

	"github.com/pkg/errors"

	"dagger.io/dagger"
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

func (fs *FS) SetDirectoryIDBy(ctx context.Context, dir *dagger.Directory) error {
	id, err := dir.ID(ctx)
	if err != nil {
		return errors.Wrap(err, "get dir id failed")
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
