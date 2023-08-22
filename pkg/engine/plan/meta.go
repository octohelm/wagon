package plan

import "github.com/octohelm/wagon/pkg/ctxutil"

var MetaContext = ctxutil.New[Meta]()

type Meta struct {
	Version string
}
