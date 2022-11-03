package core

import (
	"encoding/json"
	"github.com/octohelm/wagon/pkg/engine/daggerutil"
)

type Export struct {
	daggerutil.Exporter
}

func (m *Export) OneOf() []any {
	return []any{
		&FS{},
		&Image{},
	}
}

func (m *Export) UnmarshalJSON(data []byte) error {
	for _, v := range m.OneOf() {
		if err := json.Unmarshal(data, v); err == nil {
			if e := v.(daggerutil.Exporter); e.CanExport() {
				m.Exporter = e
				break
			}
		}
	}
	return nil
}
