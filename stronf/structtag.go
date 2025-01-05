package stronf

import "strings"

// LookupTag will return the value associated with the tag and optional path.
// The tag arg is the only required argument and uses the reflect package's
// Lookup semantics. An optional path can be provided to lookup nested values.
// Nested values can themselves be maps, each key/val pair separated by a comma.
// See examples for supported formats. The bool reports if the value was
// explicitly found at the struct tag path.
func (f Field) LookupTag(tag string, path ...string) (string, bool) {
	value, ok := f.rStructField.Tag.Lookup(tag)
	if !ok {
		return "", false
	}

	if len(path) == 0 {
		return value, true
	}

	parseStructTag := func(tag string) map[string]string {
		subTags := make(map[string]string)
		for _, subTag := range strings.Split(tag, ",") {
			key, val, _ := strings.Cut(subTag, ":")
			subTags[key] = val
		}

		return subTags
	}

	val, ok := "", false
	for _, p := range path {
		tags := parseStructTag(value)
		val, ok = tags[p]
		if !ok {
			return "", false
		}
	}

	return val, ok
}
