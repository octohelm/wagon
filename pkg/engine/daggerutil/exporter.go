package daggerutil

import "golang.org/x/net/context"

type Exporter interface {
	CanExport() bool
	ExportTo(ctx context.Context, localPath string) error
}
