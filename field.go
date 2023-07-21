package structconf

import (
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
