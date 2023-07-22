package structconf

import (
	"context"
)

type HandlerFunc func(context.Context, Field, any) (any, error)

func (h HandlerFunc) Handle(ctx context.Context, f Field, interimValue any) (any, error) {
	return h(ctx, f, interimValue)
}

var _ Handler = (HandlerFunc)(nil)

type Handler interface {
	Handle(context.Context, Field, any) (any, error)
}
