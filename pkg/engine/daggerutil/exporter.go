package daggerutil

import "golang.org/x/net/context"

type Exporter interface {
	Type() string
	CanExport() bool
	ExportTo(ctx context.Context, localPath string) error
}
