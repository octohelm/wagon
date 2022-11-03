package task

import (
	"encoding/json"
	"fmt"
	"github.com/opencontainers/go-digest"
	"github.com/spf13/afero"
	"github.com/spf13/fsync"
	"golang.org/x/mod/sumdb/dirhash"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"time"
)

func init() {
	Register(&FS{})
}

type FS struct {
	Dir  string `json:"$dir" default:""`
	Hash string `json:"$hash" default:""`
	ID   string `json:"$id" default:""`
}

func (fs *FS) Do(ctx context.Context, dgst digest.Digest) (Omit, error) {
	return nil, nil
}

func (fs *FS) HasCache() (string, bool) {
	if _, err := os.Stat(fs.Dir); err != nil {
		return "", false
	}
	return fs.Dir, true
}

func (fs *FS) Sum() error {
	hash, err := dirhash.HashDir(fs.Dir, "$wagon", dirhash.Hash1)
	if err != nil {
		return err
	}
	fs.Hash = hash
	return nil
}

func (fs *FS) DigestInputs(inputs any) (digest.Digest, error) {
	data, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}
	return digest.FromBytes(data), nil
}

func (fs *FS) NewSyncer(ctx context.Context, sourceRoot string, destRoot string) (*Syncer, error) {
	rootFs := afero.NewOsFs()
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("$wagon-cache-%d", time.Now().Unix()))
	tempDestRoot := filepath.Join(tempDir, destRoot)

	if err := os.MkdirAll(tempDestRoot, 0755); err != nil {
		return nil, err
	}

	return &Syncer{
		TempDir: tempDir,
		Dir:     fs.Dir,
		Syncer: fsync.Syncer{
			SrcFs: afero.NewReadOnlyFs(
				afero.NewBasePathFs(rootFs, sourceRoot),
			),
			DestFs: afero.NewBasePathFs(rootFs, tempDestRoot),
		},
	}, nil

}

type Syncer struct {
	fsync.Syncer
	TempDir string
	Dir     string
}

func (s *Syncer) Commit(ctx context.Context) error {
	if err := os.RemoveAll(s.Dir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.Dir), 0755); err != nil {
		return err
	}
	return os.Rename(s.TempDir, s.Dir)
}
