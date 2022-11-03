package fsutil

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	"golang.org/x/net/context"
)

type FileSyncer interface {
	Copy(ctx context.Context, src string, dst string) error
}

func NewFileSyncer(srcFs afero.Fs, dstFs afero.Fs) (FileSyncer, error) {
	if _, err := dstFs.Stat("."); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err := dstFs.MkdirAll(".", 0755); err != nil {
			return nil, err
		}
	}

	return &syncer{
		s: fsync.Syncer{
			SrcFs:  srcFs,
			DestFs: dstFs,
		},
	}, nil

}

type syncer struct {
	s fsync.Syncer
}

func (s *syncer) Copy(ctx context.Context, src string, dst string) error {
	destDir := filepath.Dir(dst)
	if destDir != "" {
		if _, err := s.s.DestFs.Stat(destDir); err != nil {
			if os.IsNotExist(err) {
				_ = s.s.DestFs.MkdirAll(destDir, 0755)
			}
		}
	}
	return s.s.Sync(dst, src)
}
