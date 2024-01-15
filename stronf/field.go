package stronf

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Field represents a single settable struct field in the parsed struct.
type Field struct {
	rVal            reflect.Value
	rStructField    reflect.StructField
	unmarshalerFunc func([]byte) error
}

// Name returns the name of the struct field.
func (f Field) Name() string {
	return f.rStructField.Name
}

// Value returns the value of the struct field.
func (f Field) Value() any {
	return f.rVal.Interface()
}

// IsZero checks if the struct field's value is it's zero type.
func (f Field) IsZero() bool {
	return f.rVal.IsZero()
}

// Kind returns the field's [reflect.Kind].
func (f Field) Kind() reflect.Kind {
	return f.rVal.Kind()
}

// Type returns the field's [reflect.Type].
func (f Field) Type() reflect.Type {
	return f.rVal.Type()
}

func (f Field) set(val any) error {
	val, err := Coerce(f, val)
	if err != nil {
		return err
	}

	if f.unmarshalerFunc != nil {
		data, ok := val.([]byte)
		if !ok {
			return fmt.Errorf("structconf: unmarshalable field %q must be provided a byte slice, got %T", f.Name(), val)
		}

		if err := f.unmarshalerFunc(data); err != nil {
			return err
		}

		return nil
	}

	rVal := reflect.ValueOf(val)

	if f.Kind() != rVal.Kind() {
		return fmt.Errorf("structconf: type mismatch, expected %q, got %q for field %q", f.Kind(), rVal.Kind(), f.Name())
	}

	f.rVal.Set(rVal)
	return nil
}

// LookupTag will return the value associated with the key and optional tag. See
// examples for supported formats. The bool reports if the key or tag was
// explicitly found in the struct tag.
func (f Field) LookupTag(key, tag string) (string, bool) {
	value, ok := f.rStructField.Tag.Lookup(key)
	if !ok {
		return "", false
	}

	if tag == "" {
		return value, true
	}

	tags := parseStructTag(value)
	val, ok := tags[tag]
	return val, ok
}

// Parse will call the handler against the field and set the field to the
// returned handler value. If the handler returns nil, no change is made to the
// field.
func (f Field) Parse(ctx context.Context, handler HandleFunc) error {
	if handler == nil {
		return errors.New("structconf: nil handler")
	}

	val, err := handler(ctx, f, nil)
	if err != nil {
		return err
	}

	if val == nil {
		return nil
	}

	if err := f.set(val); err != nil {
		return err
	}

	return nil
}
