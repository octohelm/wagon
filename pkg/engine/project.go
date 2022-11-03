package engine

import (
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/wagon/cuepkg"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	_ "github.com/octohelm/wagon/pkg/engine/plan/task"
)

func init() {
	if err := cuepkg.RegistryCueStdlibs(); err != nil {
		panic(err)
	}
}

type Project interface {
	Run(ctx context.Context, action ...string) error
}

type option struct {
	workdir   string
	entryFile string
}

type OptFunc = func(o *option)

func WithPlan(root string) OptFunc {
	return func(o *option) {
		o.entryFile = root
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

	buildConfig := cuemod.ContextFor(c.opt.workdir).BuildConfig(context.Background())

	instances := cueload.Instances([]string{c.opt.entryFile}, buildConfig)
	if len(instances) != 1 {
		return nil, errors.New("only one package is supported at a time")
	}
	c.instance = instances[0]

	return c, nil
}

type project struct {
	opt      option
	instance *build.Instance
}

func (c *project) Run(ctx context.Context, action ...string) error {
	val := cuecontext.New().BuildInstance(c.instance)
	if err := val.Err(); err != nil {
		return err
	}

	cueValue := plan.WrapValue(val)

	runner := plan.NewRunner(
		plan.NewWorkdir(c.opt.workdir, ""),
		cueValue,
	)

	return runner.Run(
		plan.TaskRunnerFactoryContext.Inject(ctx, core.DefaultFactory),
		action,
	)
}
