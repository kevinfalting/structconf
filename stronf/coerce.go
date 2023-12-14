package stronf

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Coerce will attempt to convert the provided value into the field's value.
func Coerce(field Field, val any) (any, error) {
	rVal := reflect.ValueOf(val)
	if field.unmarshalerFunc != nil {
		if !rVal.CanConvert(reflect.SliceOf(reflect.TypeOf(byte(0)))) {
			return nil, fmt.Errorf("structconf: cannot convert %q to []byte", rVal.Kind())
		}

		return rVal.Convert(reflect.SliceOf(reflect.TypeOf(byte(0)))).Interface(), nil
	}

	if rVal.Kind() == reflect.String {
		return coerceString(field, val.(string))
	}

	return val, nil
}

func coerceString(field Field, s string) (any, error) {
	switch field.Kind() {
	case reflect.String:
		return s, nil

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil

	case reflect.Float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		return float32(f), nil

	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return f, nil

	case reflect.Int:
		i, err := strconv.ParseInt(s, 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		return int(i), nil

	case reflect.Int8:
		i, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}
		return int8(i), nil

	case reflect.Int16:
		i, err := strconv.ParseInt(s, 10, 16)
		if err != nil {
			return nil, err
		}
		return int16(i), nil

	case reflect.Int32:
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(i), nil

	case reflect.Int64:
		if field.Type().String() == "time.Duration" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return nil, err
			}
			return d, nil
		}

		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return i, nil

	case reflect.Uint:
		u, err := strconv.ParseUint(s, 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		return uint(u), nil

	case reflect.Uint8:
		u, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			return nil, err
		}
		return uint8(u), nil

	case reflect.Uint16:
		u, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, err
		}
		return uint16(u), nil

	case reflect.Uint32:
		u, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(u), nil

	case reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return u, nil

	case reflect.Uintptr:
		u, err := strconv.ParseUint(s, 10, strconv.IntSize)
		if err != nil {
			return nil, err
		}
		return uintptr(u), nil

	default:
		return nil, fmt.Errorf("structconf: unsupported type: %q", field.Kind().String())
	}
}
