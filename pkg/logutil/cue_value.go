package logutil

import (
	"encoding/json"

	cueformat "cuelang.org/go/cue/format"
	"golang.org/x/exp/slog"
)

func CueValue(v any) slog.LogValuer {
	return &cueValue{v: v}
}

type cueValue struct {
	v any
}

func (c *cueValue) LogValue() slog.Value {
	data, _ := json.MarshalIndent(c.v, "", "  ")
	data, _ = cueformat.Source(data, cueformat.Simplify())
	return slog.AnyValue(data)
}
