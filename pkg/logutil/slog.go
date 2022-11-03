package logutil

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-courier/logr"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

func fromLogrLevel(l logr.Level) slog.Level {
	switch l {
	case logr.PanicLevel, logr.FatalLevel, logr.ErrorLevel:
		return slog.ErrorLevel
	case logr.WarnLevel:
		return slog.WarnLevel
	case logr.InfoLevel:
		return slog.InfoLevel
	case logr.DebugLevel:
		return slog.DebugLevel
	}
	return slog.DebugLevel
}

type slogHandler struct {
	lvl   slog.Level
	group string
	attrs []slog.Attr
}

func (s *slogHandler) Enabled(level slog.Level) bool {
	return level >= s.lvl
}

func (s *slogHandler) Handle(r slog.Record) error {
	buf := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(buf, color.WhiteString("%s", r.Time.Format("15:04:05.000")))

	switch r.Level {
	case slog.DebugLevel:
		_, _ = fmt.Fprintf(buf, " %s", color.WhiteString(strings.ToUpper(r.Level.String())[0:4]))
	case slog.WarnLevel:
		_, _ = fmt.Fprintf(buf, " %s", color.YellowString(strings.ToUpper(r.Level.String())[0:4]))
	case slog.InfoLevel:
		_, _ = fmt.Fprintf(buf, " %s", color.GreenString(strings.ToUpper(r.Level.String())[0:4]))
	case slog.ErrorLevel:
		_, _ = fmt.Fprintf(buf, " %s", color.RedString(strings.ToUpper(r.Level.String())[0:4]))
	}

	for _, attr := range s.attrs {
		if attr.Key == "name" {
			_, _ = fmt.Fprintf(buf, " %s", nameColorString(attr.Value.String()))
			break
		}
	}

	_, _ = fmt.Fprintf(buf, " %s", color.WhiteString(r.Message))

	for _, attr := range s.attrs {
		if attr.Key != "name" {
			_, _ = fmt.Fprintf(buf, " %s=%q", attr.Key, attr.Value)
		}
	}

	_, _ = fmt.Fprintf(buf, "\n")

	r.Attrs(func(attr slog.Attr) {
		if attr.Key == "err" {
			if err := attr.Value.Any().(error); err != nil {
				fmt.Println(err)
			}
		}
	})

	_, _ = io.Copy(os.Stdout, buf)

	return nil
}

func (s slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{
		lvl:   s.lvl,
		group: s.group,
		attrs: append(s.attrs, attrs...),
	}
}

func (s slogHandler) WithGroup(group string) slog.Handler {
	return &slogHandler{
		lvl:   s.lvl,
		attrs: s.attrs,
		group: group,
	}
}

func nameColorString(name string) string {
	idx, ok := colorIndexes.Load(name)
	if !ok {
		i := atomic.LoadUint32(&colorIdx)
		colorIndexes.Store(name, i)
		atomic.AddUint32(&colorIdx, 1)
		idx = i
	}
	return colorFns[int(idx.(uint32))%len(colorFns)](name)
}

type colorFn = func(fmt string, args ...any) string

var colorIndexes = sync.Map{}
var colorIdx uint32 = 0
var colorFns = []colorFn{
	color.BlueString,
	color.MagentaString,
	color.CyanString,
	color.YellowString,
}
