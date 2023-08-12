package stronf

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Field represents a single settable field in the provided struct.
type Field struct {
	reflectValue       reflect.Value
	reflectStructField reflect.StructField
}

// Name returns the name of the struct field.
func (f Field) Name() string {
	return f.reflectStructField.Name
}

// Value returns the value of the struct field.
func (f Field) Value() any {
	return f.reflectValue.Interface()
}

// IsZero checks if the struct field's value is it's zero type.
func (f Field) IsZero() bool {
	return f.reflectValue.IsZero()
}

// Zero returns the type's zero value.
func (f Field) Zero() any {
	return reflect.Zero(f.reflectValue.Type()).Interface()
}

// Kind returns the reflect package's Kind.
func (f Field) Kind() reflect.Kind {
	return f.reflectValue.Kind()
}

// Type returns the reflect package's Type.
func (f Field) Type() reflect.Type {
	return f.reflectValue.Type()
}

func (f Field) set(val any) {
	f.reflectValue.Set(val.(reflect.Value))
}

// LookupTag will return the value associated with the key and optional tag.
// Supported formats include:
// `key:"tag:val,tag1:val1"`
// `key:"tag:val", key1:"tag:val"`
// `key:"tag"` -> LookupTag("key", "") -> "tag", true
// `key:""` -> LookupTag("key", "") -> "", true
func (f Field) LookupTag(key, tag string) (string, bool) {
	value, ok := f.reflectStructField.Tag.Lookup(key)
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

func parseStructTag(tag string) map[string]string {
	options := make(map[string]string)
	for _, option := range strings.Split(tag, ",") {
		optionParts := strings.SplitN(option, ":", 2)
		key := optionParts[0]
		var value string
		if len(optionParts) > 1 {
			value = optionParts[1]
		}

		options[key] = value
	}

	return options
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

func Parse(ctx context.Context, cfg any, h Handler) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return errors.New("cfg must be a pointer")
	}

	if h == nil {
		return errors.New("handler must not be nil")
	}

	fields, err := SettableFields(cfg)
	if err != nil {
		return err
	}

	for _, field := range fields {
		val, err := h.Handle(ctx, field, nil)
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
