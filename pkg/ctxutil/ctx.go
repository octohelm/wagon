package ctxutil

import (
	contextx "github.com/octohelm/x/context"
	"golang.org/x/net/context"
)

func New[T any]() Context[T] {
	return ctx[T]{}
}

type Context[T any] interface {
	Inject(ctx context.Context, value T) context.Context
	From(ctx context.Context) T
}

type ctx[T any] struct {
}

func (c ctx[T]) Inject(ctx context.Context, value T) context.Context {
	return contextx.WithValue(ctx, c, value)
}

func (c ctx[T]) From(ctx context.Context) T {
	return ctx.Value(c).(T)
}
