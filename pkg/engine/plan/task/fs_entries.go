package task

import (
	"context"
	"dagger.io/dagger"
	"path/filepath"
	"strings"

	"github.com/octohelm/wagon/pkg/engine/daggerutil"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
)

func init() {
	core.DefaultFactory.Register(&Entries{})
}

type Entries struct {
	core.Task

	Input core.FS `json:"input"`

	Depth int `json:"depth" default:"-1"`

	// list of files and directories (ends with /) at the given path.
	Output []string `json:"-" wagon:"generated,name=output"`
}

func (e *Entries) Do(ctx context.Context) error {
	return daggerutil.Do(ctx, func(c *dagger.Client) error {
		fw := &entriesWalker{
			d: c.Directory(dagger.DirectoryOpts{
				ID: e.Input.DirectoryID(),
			}),
			maxDepth: e.Depth,
		}

		if err := fw.walk(ctx, "/", func(path string) error {
			e.Output = append(e.Output, path)
			return nil
		}); err != nil {
			return err
		}

		return nil
	})
}

type entriesWalker struct {
	d        *dagger.Directory
	maxDepth int
}

func (w *entriesWalker) walk(ctx context.Context, path string, walkFunc func(path string) error) error {
	if path == "" {
		path = "/"
	}

	entries, err := w.d.Entries(ctx, dagger.DirectoryEntriesOpts{
		Path: path,
	})
	if err != nil {
		// entries failed means path is not directory
		if err := walkFunc(path); err != nil {
			return err
		}
		return nil
	}

	if w.maxDepth > 0 && strings.Count(path, "/") >= w.maxDepth {
		return walkFunc(path + "/")
	}

	if len(entries) == 0 {
		return walkFunc(path + "/")
	}

	for _, entry := range entries {
		p := filepath.Join(path, entry)

		if err := w.walk(ctx, p, walkFunc); err != nil {
			return err
		}
	}

	return nil
}
