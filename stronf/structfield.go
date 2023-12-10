package stronf

import (
	"encoding"
	"errors"
	"reflect"
)

// SettableFields returns a slice of all settable struct fields in the provided
// struct. The provided argument must be a pointer to a struct.
func SettableFields(v any) ([]Field, error) {
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Pointer {
		return nil, errors.New("structconf: must be pointer")
	}

	rVal = rVal.Elem()
	if rVal.Kind() != reflect.Struct {
		return nil, errors.New("structconf: must be a struct")
	}

	var fields []Field
	if err := settableFields(rVal, &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

func settableFields(rVal reflect.Value, fields *[]Field) error {
	for i := 0; i < rVal.NumField(); i++ {
		rValField := rVal.Field(i)
		rStructField := rVal.Type().Field(i)

		if !rValField.CanSet() {
			continue
		}

		unmarshaler := unmarshalerFunc(rValField)
		if unmarshaler != nil {
			*fields = append(*fields, Field{
				rVal:            rValField,
				rStructField:    rStructField,
				unmarshalerFunc: unmarshaler,
			})

			continue
		}

		switch rValField.Kind() {
		case reflect.Pointer:
			continue

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Complex64, reflect.Complex128,
			reflect.Float32, reflect.Float64,
			reflect.Uintptr,
			reflect.String,
			reflect.Bool:
			*fields = append(*fields, Field{
				rVal:         rValField,
				rStructField: rStructField,
			})

		case reflect.Struct:
			if err := settableFields(rValField, fields); err != nil {
				return err
			}

			continue

		default:
			// unsupported kind
			continue
		}
	}

	return nil
}

func unmarshalerFunc(rVal reflect.Value) func([]byte) error {
	if !rVal.CanAddr() {
		return nil
	}

	if textUnmarshaler, ok := rVal.Addr().Interface().(encoding.TextUnmarshaler); ok {
		return textUnmarshaler.UnmarshalText
	}

	if binaryUnmarshaler, ok := rVal.Addr().Interface().(encoding.BinaryUnmarshaler); ok {
		return binaryUnmarshaler.UnmarshalBinary
	}

	return nil
}
