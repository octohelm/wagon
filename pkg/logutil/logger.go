package logutil

import (
	"github.com/go-courier/logr"
	"github.com/innoai-tech/infra/pkg/configuration"
	"github.com/innoai-tech/infra/pkg/otel"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
)

type Logger struct {
	// Log level
	LogLevel otel.LogLevel `flag:",omitempty"`
	logger   logr.Logger
}

func (o *Logger) SetDefaults() {
	if o.LogLevel == "" {
		o.LogLevel = otel.InfoLevel
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
