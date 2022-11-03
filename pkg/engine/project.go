package engine

import (
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
	"github.com/go-courier/logr"
	"github.com/octohelm/cuemod/pkg/cuemod"
	"github.com/octohelm/wagon/cuepkg"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"github.com/octohelm/wagon/pkg/engine/plan/task/core"
	"github.com/octohelm/wagon/pkg/fsutil"
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
	output    string
}

type OptFunc = func(o *option)

func WithPlan(root string) OptFunc {
	return func(o *option) {
		o.entryFile = root
	}
}

func WithOutput(output string) OptFunc {
	return func(o *option) {
		o.output = output
	}
}

func WithWorkdir(workdir string) OptFunc {
	return func(o *option) {
		o.workdir = workdir
	}
}

func New(ctx context.Context, opts ...OptFunc) (Project, error) {
	c := &project{}
	for i := range opts {
		opts[i](&c.opt)
	}

	buildConfig := cuemod.ContextFor(c.opt.workdir).BuildConfig(ctx)

	instances := cueload.Instances([]string{c.opt.entryFile}, buildConfig)
	if len(instances) != 1 {
		return nil, errors.New("only one package is supported at a time")
	}
	c.instance = instances[0]

	if err := c.instance.Err; err != nil {
		return nil, err
	}

	return c, nil
}

type project struct {
	opt      option
	instance *build.Instance
}

func (c *project) Run(ctx context.Context, action ...string) error {
	logr.FromContext(ctx).WithValues("name", "Project").Debug("loading...")

	val := cuecontext.New().BuildInstance(c.instance)
	if err := val.Err(); err != nil {
		return err
	}

	cueValue := plan.WrapValue(val)
	workdir := plan.NewWorkdir(c.opt.workdir, "")
	registryAuthStore := plan.NewRegistryAuthStore()

	logr.FromContext(ctx).WithValues("name", "Project").Debug("starting...")

	runner := plan.NewRunner(cueValue)

	ctx = plan.TaskRunnerFactoryContext.Inject(ctx, core.DefaultFactory)
	ctx = plan.WorkdirContext.Inject(ctx, workdir)
	ctx = plan.RegistryAuthStoreContext.Inject(ctx, registryAuthStore)

	output, err := runner.Run(ctx, action)
	if err != nil {
		return err
	}

	if output != nil && c.opt.output != "" {
		fs := workdir.Fs(plan.WorkdirProject, c.opt.output)
		dest, err := fsutil.RealPath(fs)
		if err != nil {
			return err
		}

		e := &core.Export{}
		if err := output.Decode(e); err != nil {
			return err
		}

		if e.Exporter != nil {
			l := logr.FromContext(ctx).WithValues(
				"name", "Export",
				"dest", c.opt.output,
			)
			return e.Exporter.ExportTo(logr.WithLogger(ctx, l), dest)
		}
	}

	return nil
}
