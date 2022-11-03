package plan

import (
	"github.com/octohelm/wagon/pkg/engine/spec"
)

type Context interface {
	Workdir() Workdir
	Pkg() spec.Pkg
	Target() spec.Platform
}
