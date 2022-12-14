package logutil

import (
	"context"
	"fmt"
	"github.com/go-courier/logr"
	"golang.org/x/exp/slog"
	"strings"
)

type logger struct {
	slogr slog.Logger
	spans []string
}

func (d *logger) SetLevel(lvl logr.Level) {
	d.slogr.Handler().(*slogHandler).lvl = fromLogrLevel(lvl)
}

func (d *logger) WithValues(keyAndValues ...any) logr.Logger {
	return &logger{
		spans: d.spans,
		slogr: d.slogr.With(keyAndValues...),
	}
}

func (d *logger) Start(ctx context.Context, name string, keyAndValues ...any) (context.Context, logr.Logger) {
	spans := append(d.spans, name)

	return ctx, &logger{
		spans: spans,
		slogr: d.slogr.WithGroup(strings.Join(spans, "/")).With(keyAndValues...),
	}
}

func (d *logger) End() {
	if len(d.spans) != 0 {
		d.spans = d.spans[0 : len(d.spans)-1]
	}
}

func (d *logger) Debug(format string, args ...any) {
	d.slogr.Debug(fmt.Sprintf(format, args...))
}

func (d *logger) Info(format string, args ...any) {
	d.slogr.Info(fmt.Sprintf(format, args...))
}

func (d *logger) Warn(err error) {
	d.slogr.Warn(err.Error())
}

func (d *logger) Error(err error) {
	d.slogr.Error(err.Error(), err)
}

func (d *logger) Trace(format string, args ...any) {
	d.Debug(fmt.Sprintf(format, args...))
}

func (d *logger) Fatal(err error) {
	d.Error(err)
}

func (d *logger) Panic(err error) {
	d.Error(err)
}
