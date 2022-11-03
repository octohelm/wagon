package old

import (
	"github.com/octohelm/wagon/pkg/engine/plan"
	"golang.org/x/net/context"
)

func init() {
	Register("Context", &targetPlatformTask{})
}

type targetPlatformTask struct {
	*plan.Task
}

func (targetPlatformTask) New(task *plan.Task) plan.TaskRunner {
	return &targetPlatformTask{Task: task}
}

func (targetPlatformTask) Decl() []byte {
	return []byte(`
workdir: {
	pwd: string
	source: string
	output: string
}
target: {
	os:           string
	arch:         string
	variant?:     string
	osVersion?:   string
	osFeatures?:  [...string]
}
`)
}

func (t *targetPlatformTask) Run(ctx context.Context, planCtx plan.Context) error {
	p := planCtx.Target()
	workdir := planCtx.Workdir()

	return t.Fill(map[string]any{
		"workdir": map[string]any{
			"pwd":    workdir.Pwd(),
			"source": workdir.Source(),
			"output": workdir.Output(),
		},
		"target": map[string]any{
			"os":         p.OS,
			"arch":       p.Arch,
			"variant":    p.Variant,
			"osVersion":  p.OSVersion,
			"osFeatures": p.OSFeatures,
		},
	})
}
