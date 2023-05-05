package logutil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/mattn/go-colorable"

	"golang.org/x/net/context"

	"github.com/go-courier/logr"
	"golang.org/x/exp/slog"
)

func fromLogrLevel(l logr.Level) slog.Level {
	switch l {
	case logr.ErrorLevel:
		return slog.LevelError
	case logr.WarnLevel:
		return slog.LevelWarn
	case logr.InfoLevel:
		return slog.LevelInfo
	case logr.DebugLevel:
		return slog.LevelDebug
	}
	return slog.LevelDebug
}

type slogHandler struct {
	lvl   slog.Level
	group string
	attrs []slog.Attr
}

func (s *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= s.lvl
}

type Printer interface {
	Fprintf(w io.Writer, format string, a ...any) (n int, err error)
}

func withLevelColor(level slog.Level) func(io.Writer) io.Writer {
	switch level {
	case slog.LevelError:
		return WithColor(FgRed)
	case slog.LevelWarn:
		return WithColor(FgYellow)
	case slog.LevelInfo:
		return WithColor(FgGreen)
	}
	return WithColor(FgWhite)
}

func (s *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	prefix := bytes.NewBuffer(nil)

	_, _ = fmt.Fprintf(WithColor(FgWhite)(prefix), "%s", r.Time.Format("15:04:05.000"))
	_, _ = fmt.Fprintf(withLevelColor(r.Level)(prefix), " %s", strings.ToUpper(r.Level.String())[0:4])

	for _, attr := range s.attrs {
		if attr.Key == "name" {
			_, _ = fmt.Fprintf(WithNameColor(attr.Value.String())(prefix), " %s |", attr.Value.String())
			break
		}
	}

	w := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(bytes.NewBufferString(r.Message))
	idx := 0

	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); len(line) > 0 {
			if idx > 0 {
				_, _ = fmt.Fprintln(w)
			}
			_, _ = fmt.Fprintf(w, "%s %s", prefix.String(), line)
			idx++
		}
	}

	for _, attr := range s.attrs {
		if attr.Key != "name" {
			switch attr.Value.Kind() {
			case slog.KindString:
				_, _ = fmt.Fprintf(WithColor(FgWhite)(w), " %s=%q", attr.Key, attr.Value)
			default:
				logValue := attr.Value.Any()
				if valuer, ok := logValue.(slog.LogValuer); ok {
					logValue = valuer.LogValue().Any()
				}

				switch x := logValue.(type) {
				case []byte:
					_, _ = fmt.Fprintf(WithColor(FgWhite)(w), " %s=%v", attr.Key, string(x))
				default:
					_, _ = fmt.Fprintf(WithColor(FgWhite)(w), " %s=%v", attr.Key, x)
				}
			}
		}
	}

	_, _ = fmt.Fprintln(w)

	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "err" {
			if err := attr.Value.Any().(error); err != nil {
				_, _ = fmt.Fprintf(w, "%+v", err)
			}
		}
		return true
	})

	_, _ = io.Copy(colorable.NewColorableStdout(), w)

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

var colorIndexes = sync.Map{}
var colorIdx uint32 = 0
var colorFns = []WrapWriter{
	WithColor(FgBlue),
	WithColor(FgMagenta),
	WithColor(FgCyan),
	WithColor(FgYellow),
}

func WithNameColor(name string) WrapWriter {
	idx, ok := colorIndexes.Load(name)
	if !ok {
		i := atomic.LoadUint32(&colorIdx)
		colorIndexes.Store(name, i)
		atomic.AddUint32(&colorIdx, 1)
		idx = i
	}
	return colorFns[int(idx.(uint32))%len(colorFns)]
}
