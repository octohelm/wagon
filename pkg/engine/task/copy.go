package task

import (
	"fmt"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/fsutil"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
	"path/filepath"
)

func init() {
	Register(&Copy{})
}

type Copy struct {
	Input    FS       `json:"input"`
	Contents FS       `json:"contents"`
	Source   string   `json:"source" default:"/"`
	Dest     string   `json:"dest" default:"/"`
	Include  []string `json:"include" default:"[\"*\"]"`
	Exclude  []string `json:"exclude"`

	FS `json:"output,omitempty"`
}

func (c Copy) Do(ctx context.Context, dgst digest.Digest) (Omit, error) {
	workdir := plan.WorkdirContext.From(ctx)
	log := logr.FromContext(ctx)

	c.ID = string(dgst)
	c.Dir = workdir.Cache(c.ID)

	_, ok := c.HasCache()
	if !ok {
		f, err := fsutil.NewFinder(c.Include, c.Exclude)
		if err != nil {
			return nil, err
		}

		root := filepath.Join(c.Contents.Dir, c.Source)

		w, err := c.NewSyncer(ctx, root, c.Dest)
		if err != nil {
			return nil, err
		}

		if err := f.Walk(root, func(filename string) error {
			log.Debug(fmt.Sprintf("COPY %s", filename))
			return w.Sync(filename, filename)
		}); err != nil {
			return nil, err
		}

		if err := w.Commit(ctx); err != nil {
			return nil, err
		}
	}

	return func(t plan.Task) error {
		return t.Fill(map[string]any{
			"output": c.FS,
		})
	}, nil
}
