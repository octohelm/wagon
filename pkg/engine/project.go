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
	"github.com/octohelm/wagon/pkg/version/gomod"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"os"
	"path"
	"strings"

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
	entryFile       string
	output          string
	imagePullPrefix string
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

var inCI = os.Getenv("CI") == "true"

func New(ctx context.Context, opts ...OptFunc) (Project, error) {
	c := &project{}
	for i := range opts {
		opts[i](&c.opt)
	}

	cwd, _ := os.Getwd()
	sourceRoot := path.Join(cwd, c.opt.entryFile)

	if strings.HasSuffix(sourceRoot, ".cue") {
		sourceRoot = path.Dir(sourceRoot)
	}

	c.sourceRoot = sourceRoot

	buildConfig := cuemod.ContextFor(cwd).BuildConfig(ctx)

	instances := cueload.Instances([]string{c.opt.entryFile}, buildConfig)
	if len(instances) != 1 {
		return nil, errors.New("only one package is supported at a time")
	}
	c.instance = instances[0]

	if err := c.instance.Err; err != nil {
		return nil, err
	}

	version, err := gomod.LocalRevInfo(cwd)
	if err != nil {
		return nil, errors.Wrap(err, "local rev failed")
	}

	if inCI && strings.Contains(version.Version, "-dirty") {
		return nil, errors.New("dirty build not allowed in CI")
	}

	c.Version = version.Version

	return c, nil
}

type project struct {
	opt        option
	sourceRoot string
	instance   *build.Instance
	plan.Meta
}

func (c *project) Run(ctx context.Context, action ...string) error {
	val := cuecontext.New().BuildInstance(c.instance)
	if err := val.Err(); err != nil {
		return err
	}

	cueValue := plan.WrapValue(val)
	workdir := plan.NewWorkdir(c.sourceRoot, "")
	registryAuthStore := plan.NewRegistryAuthStore()

	l := logr.FromContext(ctx)

	runner := plan.NewRunner(cueValue, c.opt.output, &core.FS{}, &core.Image{})

	ctx = plan.TaskRunnerFactoryContext.Inject(ctx, core.DefaultFactory)
	ctx = plan.WorkdirContext.Inject(ctx, workdir)
	ctx = plan.RegistryAuthStoreContext.Inject(ctx, registryAuthStore)
	ctx = plan.MetaContext.Inject(ctx, c.Meta)
	ctx = logr.WithLogger(ctx, l.WithValues("version", c.Version))

	return runner.Run(ctx, action)
}
