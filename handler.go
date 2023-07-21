package structconf

import (
	"context"
)

type HandlerFunc func(context.Context, Field) (any, error)

func (h HandlerFunc) Handle(ctx context.Context, f Field) (any, error) {
	return h(ctx, f)
}

var _ Handler = (HandlerFunc)(nil)

type Handler interface {
	Handle(context.Context, Field) (any, error)
}
