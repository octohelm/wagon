package daggerutil

import (
	"dagger.io/dagger"
	"golang.org/x/net/context"
)

func Query(ctx context.Context, c *dagger.Client, data interface{}, query string) error {
	return c.Do(
		ctx,
		&dagger.Request{
			Query: query,
		},
		&dagger.Response{
			Data: data,
		},
	)
}
