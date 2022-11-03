package task

import (
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/opencontainers/go-digest"
	"golang.org/x/net/context"
)

func init() {
	Register(&Source{})
}

type Source struct {
	Path string `json:"path" default:"."`

	FS `json:"output"`
}

func (s *Source) PreDo(ctx context.Context) error {
	w := plan.WorkdirContext.From(ctx)
	s.Dir = w.Source(s.Path)
	return s.Sum()
}

func (s Source) Do(ctx context.Context, dgst digest.Digest) (Omit, error) {
	return func(t plan.Task) error {
		return t.Fill(map[string]any{
			"output": s.FS,
		})
	}, nil
}
