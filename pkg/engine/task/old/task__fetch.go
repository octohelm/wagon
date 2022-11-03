package old

import (
	"fmt"
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/fsutil"
	"net/url"

	cueerrors "cuelang.org/go/cue/errors"
	"github.com/octohelm/wagon/pkg/downloader"
	"github.com/octohelm/wagon/pkg/engine/plan"
	"golang.org/x/net/context"
)

func init() {
	//Register("Fetch", &fetchTask{})
}

type fetchTask struct {
	*plan.Task
}

func (fetchTask) Decl() []byte {
	return []byte(`
source:    string
dest:   string
`)
}

func (fetchTask) New(task *plan.Task) plan.TaskRunner {
	return &fetchTask{Task: task}
}

func (t *fetchTask) Run(ctx context.Context, planCtx plan.Context) error {
	workdir := planCtx.Workdir()
	val := t.Value()

	var source string

	if err := val.Lookup("source").Decode(&source); err != nil {
		return err
	}

	u, err := url.Parse(source)
	if err != nil {
		return cueerrors.Newf(val.Pos(), "`source` of #Fetch must be valid url")
	}

	d := downloader.NewDownloader(workdir.Cache())

	logr.FromContext(ctx).Info(fmt.Sprintf("fetch %s", u.String()))

	info, err := d.Resolve(ctx, u)
	if err != nil {
		return cueerrors.Newf(val.Pos(), "download failed")
	}

	var output string

	_ = val.Lookup("output").Decode(&output)

	if output != "" {
		_, err := fsutil.Copy(info.File, output)
		return err
	}

	return t.Fill(map[string]any{
		"output": info.File,
	})
}
