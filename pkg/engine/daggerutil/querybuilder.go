package daggerutil

import (
	"dagger.io/dagger"
	"golang.org/x/net/context"
)

type QueryClient interface {
	Query(ctx context.Context, data interface{}, query string) error
}

func NewQueryClient(c *dagger.Client) QueryClient {
	return &queryClient{c: c}
}

type queryClient struct {
	c *dagger.Client
}

func (c *queryClient) Query(ctx context.Context, data interface{}, query string) error {
	return c.c.Do(ctx, &dagger.Request{Query: query}, &dagger.Response{Data: data})
}
