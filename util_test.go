package structconf_test

import (
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf"
)

func TestParseStringForKind(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		kind    reflect.Kind
		want    interface{}
		wantErr bool
	}{
		{name: "Parse string", input: "hello", kind: reflect.String, want: "hello"},
		{name: "Parse bool - true", input: "true", kind: reflect.Bool, want: true},
		{name: "Parse bool - false", input: "false", kind: reflect.Bool, want: false},
		{name: "Parse bool - invalid", input: "not-a-bool", kind: reflect.Bool, wantErr: true},
		{name: "Parse float32 - valid", input: "3.14", kind: reflect.Float32, want: float32(3.14)},
		{name: "Parse float32 - invalid", input: "not-a-float", kind: reflect.Float32, wantErr: true},
		{name: "Parse float64 - valid", input: "3.14", kind: reflect.Float64, want: 3.14},
		{name: "Parse float64 - invalid", input: "not-a-float", kind: reflect.Float64, wantErr: true},
		{name: "Parse int - valid", input: "100", kind: reflect.Int, want: 100},
		{name: "Parse int - invalid", input: "not-an-int", kind: reflect.Int, wantErr: true},
		{name: "Parse int8 - valid", input: "100", kind: reflect.Int8, want: int8(100)},
		{name: "Parse int8 - overflow", input: "130", kind: reflect.Int8, wantErr: true},
		{name: "Parse int16 - valid", input: "100", kind: reflect.Int16, want: int16(100)},
		{name: "Parse int32 - valid", input: "100", kind: reflect.Int32, want: int32(100)},
		{name: "Parse int64 - valid", input: "100", kind: reflect.Int64, want: int64(100)},
		{name: "Parse uint - valid", input: "100", kind: reflect.Uint, want: uint(100)},
		{name: "Parse uint - invalid", input: "-1", kind: reflect.Uint, wantErr: true},
		{name: "Parse uint8 - valid", input: "100", kind: reflect.Uint8, want: uint8(100)},
		{name: "Parse uint8 - overflow", input: "300", kind: reflect.Uint8, wantErr: true},
		{name: "Parse uint16 - valid", input: "100", kind: reflect.Uint16, want: uint16(100)},
		{name: "Parse uint32 - valid", input: "100", kind: reflect.Uint32, want: uint32(100)},
		{name: "Parse uint64 - valid", input: "100", kind: reflect.Uint64, want: uint64(100)},
		{name: "Parse uintptr - valid", input: "100", kind: reflect.Uintptr, want: uintptr(100)},
		{name: "Parse uintptr - invalid", input: "-1", kind: reflect.Uintptr, wantErr: true},
		{name: "Unsupported type", input: "100", kind: reflect.Slice, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := structconf.ParseStringForKind(tt.input, tt.kind)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
