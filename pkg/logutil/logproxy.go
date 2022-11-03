package logutil

import (
	"github.com/go-courier/logr"
	"github.com/octohelm/wagon/pkg/ctxutil"
	"sync"
)

var LogProxyContext = ctxutil.New[*LogProxy]()

type LogProxy struct {
	l  logr.Logger
	rw sync.RWMutex
}

func (p *LogProxy) Bind(l logr.Logger) func() {
	p.rw.Lock()
	defer p.rw.Unlock()

	p.l = l

	return func() {
		p.rw.Lock()
		defer p.rw.Unlock()

		p.l = nil
	}
}

func (p *LogProxy) Info(format string, args ...any) {
	p.rw.Lock()
	defer p.rw.Unlock()

	if p.l != nil {
		p.l.Info(format, args...)
	}
}
