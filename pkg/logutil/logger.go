package logutil

import (
	"github.com/go-courier/logr"
	"github.com/innoai-tech/infra/pkg/configuration"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
)

// +gengo:enum
type LogLevel string

const (
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
	InfoLevel  LogLevel = "info"
	DebugLevel LogLevel = "debug"
)

type Logger struct {
	// Log level
	LogLevel LogLevel `flag:",omitempty"`
	logger   logr.Logger
}

func (o *Logger) SetDefaults() {
	if o.LogLevel == "" {
		o.LogLevel = InfoLevel
	}
}

func (o *Logger) Init(ctx context.Context) error {
	if o.logger == nil {
		lvl, _ := logr.ParseLevel(string(o.LogLevel))
		o.logger = &logger{slogr: slog.New(&slogHandler{lvl: fromLogrLevel(lvl)})}
	}
	return nil
}

func (o *Logger) InjectContext(ctx context.Context) context.Context {
	return configuration.InjectContext(
		ctx,
		configuration.InjectContextFunc(logr.WithLogger, o.logger),
	)
}
