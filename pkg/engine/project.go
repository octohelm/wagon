package engine

import (
	_ "embed"
	"os"
	"path/filepath"

	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/cuecontext"
	astparser "cuelang.org/go/cue/parser"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/spec"
	"github.com/octohelm/wagon/pkg/engine/task"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

//go:embed types.cue
var declTypes []byte
var declFile = "wagon.cue"

type Project interface {
	Run(ctx context.Context, target spec.Platform) error
}

type option struct {
	root    string
	workdir string
}

type OptFunc = func(o *option)

func WithProjectRoot(root string) OptFunc {
	return func(o *option) {
		o.root = root
	}
}

func WithWorkdir(workdir string) OptFunc {
	return func(o *option) {
		o.workdir = workdir
	}
}

func New(opts ...OptFunc) (Project, error) {
	c := &project{}
	for i := range opts {
		opts[i](&c.opt)
	}

	f := filepath.Join(c.opt.root, declFile)

	data, err := os.ReadFile(filepath.Join(c.opt.workdir, f))
	if err != nil {
		return nil, errors.Wrapf(err, "read `%s` failed", f)
	}

	combined := append(data, append([]byte("\n"), declTypes...)...)
	combined = append(combined, task.Default.AllDefs()...)

	file, err := astparser.ParseFile(f, combined, astparser.ParseComments)
	if err != nil {
		return nil, errors.Wrap(err, "parse failed")
	}
	c.file = file

	return c, nil
}

type project struct {
	opt  option
	file *cueast.File
}

func (c *project) Run(ctx context.Context, target spec.Platform) error {
	cueValue := plan.WrapValue(cuecontext.New().BuildFile(c.file))

	pkg := spec.Pkg{}

	pkgInfoValue := cueValue.Lookup("pkg")
	if err := pkgInfoValue.Decode(&pkg); err != nil {
		return err
	}

	root := plan.WorkdirFor(pkg, target, c.opt.root, c.opt.workdir)
	runner := plan.NewRunner(pkg, root, target, cueValue)

	return runner.Run(plan.TaskRunnerResolverContext.Inject(ctx, task.Default))
}
