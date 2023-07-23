package structconf

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Conf holds the handlers and middleware for parsing structs. Handlers and
// Middleware should be provided in the order they're intended to be run, each
// taking precedence over the last.
type Conf[T any] struct {
	Handlers   []Handler
	Middleware []Middleware
}

var _ Handler = (*Conf[any])(nil)

// New returns a Conf initialized with all of the default Handlers and
// Middleware to parse a struct.
func New[T any](opts ...confOptionFunc) (*Conf[T], error) {
	var confOpt confOption
	for _, opt := range opts {
		opt(&confOpt)
	}

	flagHandler, err := NewFlag[T](confOpt.flagSet)
	if err != nil {
		return nil, err
	}

	conf := Conf[T]{
		Handlers: []Handler{
			EnvironmentVariable{},
			flagHandler,
		},

		Middleware: []Middleware{
			Required(),
			Default(),
		},
	}

	if confOpt.rsaPrivateKey != nil {
		conf.Handlers = append(conf.Handlers, &RSAHandler{
			PrivateKey: confOpt.rsaPrivateKey,
			Label:      confOpt.rsaLabel,
		})
	}

	return &conf, nil
}

// Handle will call all of the configured middleware and handlers on a single
// field.
func (c *Conf[T]) Handle(ctx context.Context, field Field, interimValue any) (any, error) {
	h := WrapMiddleware(c.Handlers, c.Middleware...)
	return h.Handle(ctx, field, interimValue)
}

// SettableFields returns a slice of all settable struct fields in the provided
// struct. The provided argument must be a pointer to a struct.
func SettableFields(i any) ([]Field, error) {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %v", val.Kind())
	}

	var fields []Field
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Pointer:
			if fieldValue.Elem().Kind() != reflect.Struct || fieldValue.IsNil() {
				continue
			}

			f, err := SettableFields(fieldValue.Interface())
			if err != nil {
				return nil, err
			}

			fields = append(fields, f...)

		case reflect.Struct:
			if fieldValue.CanAddr() || fieldType.Anonymous {
				f, err := SettableFields(fieldValue.Addr().Interface())
				if err != nil {
					return nil, err
				}

				fields = append(fields, f...)

			} else {
				f, err := SettableFields(fieldValue.Interface())
				if err != nil {
					return nil, err
				}

				fields = append(fields, f...)
			}

			continue

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Complex64, reflect.Complex128,
			reflect.Float32, reflect.Float64,
			reflect.Uintptr,
			reflect.String,
			reflect.Bool:

			if !fieldValue.CanSet() {
				continue
			}

			fields = append(fields, Field{
				reflectStructField: fieldType,
				reflectValue:       fieldValue,
			})

		case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func,
			reflect.UnsafePointer, reflect.Array, reflect.Interface:
			// Skip fields of these kinds
			continue

		default:
			return nil, fmt.Errorf("unexpected kind: %v", fieldValue.Kind())
		}
	}

	return fields, nil
}

// Parse will walk every field and nested field in the provided struct appling
// the handlers to each field.
func (c *Conf[T]) Parse(ctx context.Context, cfg *T) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("cfg must be a pointer")
	}

	fields, err := SettableFields(cfg)
	if err != nil {
		return err
	}

	for _, field := range fields {
		val, err := c.Handle(ctx, field, nil)
		if err != nil {
			return err
		}

		if val == nil {
			continue
		}

		field.set(reflect.ValueOf(val))
	}

	return nil
}
