package structconf

import (
	"context"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

// Conf holds the handlers and middleware for parsing structs. Handlers and
// Middleware should be provided in the order they're intended to be run, each
// taking precedence over the last.
type Conf[T any] struct {
	Handlers   []stronf.Handler
	Middleware []stronf.Middleware
}

var _ stronf.Handler = (*Conf[any])(nil)

// New returns a Conf initialized with all of the default Handlers and
// Middleware to parse a struct.
func New[T any](opts ...confOptionFunc) (*Conf[T], error) {
	var confOpt confOption
	for _, opt := range opts {
		opt(&confOpt)
	}

	flagHandler, err := confhandler.NewFlag[T](confOpt.flagSet)
	if err != nil {
		return nil, err
	}

	conf := Conf[T]{
		Handlers: []stronf.Handler{
			confhandler.EnvironmentVariable{},
			flagHandler,
		},

		Middleware: []stronf.Middleware{
			confhandler.Required(),
			confhandler.Default(),
		},
	}

	if confOpt.rsaPrivateKey != nil {
		conf.Handlers = append(conf.Handlers, &confhandler.RSAHandler{
			PrivateKey: confOpt.rsaPrivateKey,
			Label:      confOpt.rsaLabel,
		})
	}

	return &conf, nil
}

// Handle will call all of the configured middleware and handlers on a single
// field.
func (c *Conf[T]) Handle(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
	handler := stronf.CombineHandlers(c.Handlers...)
	handler = stronf.WrapMiddleware(handler, c.Middleware...)
	return handler.Handle(ctx, field, interimValue)
}

// Parse will walk every field and nested field in the provided struct appling
// the handlers to each field.
func (c *Conf[T]) Parse(ctx context.Context, cfg *T) error {
	fields, err := stronf.SettableFields(cfg)
	if err != nil {
		return err
	}

	for _, field := range fields {
		if err := field.Parse(ctx, c); err != nil {
			return err
		}
	}

	return nil
}
