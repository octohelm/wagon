package gomod

import (
	"github.com/octohelm/cuemod/pkg/modutil"
	"golang.org/x/net/context"
)

func LocalRevInfo(workdir string) (*modutil.RevInfo, error) {
	return modutil.RevInfoFromDir(context.Background(), workdir)
}
