package stronf_test

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
	"unsafe"

	"github.com/kevinfalting/structconf/stronf"
)

func TestLookupTag(t *testing.T) {
	t.Run("key:\"tag:val,tag1:val1\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag:val,tag1:val1"`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "tag")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val" {
			t.Errorf("expected v = val, got %q", v)
		}
	})

	t.Run("key:\"\"", func(t *testing.T) {
		A := struct {
			Int int `key:""`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "")
		if !ok {
			t.Error("expected key to exist")
		}

		if v != "" {
			t.Errorf("expected v = \"\", got %q", v)
		}
	})

	t.Run("key:\"tag\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag"`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "")
		if !ok {
			t.Error("expected key to exist")
		}

		if v != "tag" {
			t.Errorf("expected v = tag, got %q", v)
		}
	})

	t.Run("key:\"tag:val\" key1:\"tag1:val1\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag:val" key1:"tag1:val1"`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "tag")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val" {
			t.Errorf("expected v = val, got %q", v)
		}

		v, ok = field.LookupTag("key1", "tag1")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val1" {
			t.Errorf("expected v = val, got %q", v)
		}
	})
}

func TestSettableFields(t *testing.T) {
	{
		errTests := []struct {
			name string
			arg  any
		}{
			{
				name: "nil",
				arg:  nil,
			},
			{
				name: "interface",
				arg:  new(any),
			},
			{
				name: "slice",
				arg:  make([]int, 0),
			},
			{
				name: "array",
				arg:  [1]int{},
			},
			{
				name: "map",
				arg:  make(map[any]any),
			},
			{
				name: "chan",
				arg:  make(chan any),
			},
			{
				name: "boolean",
				arg:  true,
			},
			{
				name: "integer",
				arg:  int(1),
			},
			{
				name: "integer8",
				arg:  int8(1),
			},
			{
				name: "integer16",
				arg:  int16(1),
			},
			{
				name: "integer32",
				arg:  int32(1),
			},
			{
				name: "integer64",
				arg:  int64(1),
			},
			{
				name: "unsigned integer",
				arg:  uint(1),
			},
			{
				name: "unsigned integer8",
				arg:  uint8(1),
			},
			{
				name: "byte",
				arg:  byte(1),
			},
			{
				name: "unsigned integer16",
				arg:  uint16(1),
			},
			{
				name: "unsigned integer32",
				arg:  uint32(1),
			},
			{
				name: "unsigned integer64",
				arg:  uint64(1),
			},
			{
				name: "unsigned integer pointer",
				arg:  uintptr(1),
			},
			{
				name: "float32",
				arg:  float32(1.0),
			},
			{
				name: "float64",
				arg:  float64(1.0),
			},
			{
				name: "complex64",
				arg:  complex64(1 + 2i),
			},
			{
				name: "complex128",
				arg:  complex128(1 + 2i),
			},
			{
				name: "string",
				arg:  "test",
			},
			{
				name: "rune",
				arg:  rune('a'),
			},
			{
				name: "function",
				arg:  func() {},
			},
			{
				name: "unsafe pointer",
				arg:  unsafe.Pointer(new(int)),
			},
		}

		for _, tt := range errTests {
			t.Run(tt.name+" returns error", func(t *testing.T) {
				_, err := stronf.SettableFields(tt.arg)
				if err == nil {
					t.Fatalf("expected error, got %v", err)
				}
			})
		}
	}

	{
		type TestStruct struct {
			Field1 int
			Field2 string
		}

		noErrTests := []struct {
			name string
			arg  any
		}{
			{
				name: "struct",
				arg:  TestStruct{},
			},
			{
				name: "pointer to struct",
				arg:  &TestStruct{},
			},
		}

		for _, tt := range noErrTests {
			t.Run(tt.name+" does not return error", func(t *testing.T) {
				_, err := stronf.SettableFields(tt.arg)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			})
		}
	}

	{
		type TestStruct struct {
			Field1 int
			Field2 string
		}

		type testStruct struct {
			Field3 int
			Field4 string
		}

		type AllTypesStruct struct {
			// Basic types - 19 / 19
			BoolValue       bool
			IntValue        int
			Int8Value       int8
			Int16Value      int16
			Int32Value      int32
			Int64Value      int64
			UintValue       uint
			Uint8Value      uint8
			ByteValue       byte
			Uint16Value     uint16
			Uint32Value     uint32
			Uint64Value     uint64
			UintptrValue    uintptr
			Float32Value    float32
			Float64Value    float64
			Complex64Value  complex64
			Complex128Value complex128
			StringValue     string
			RuneValue       rune

			// Composite types - 6 / 0
			ArrayValue   [1]int
			SliceValue   []int
			MapValue     map[string]int
			ChanValue    chan int
			FuncValue    func() int
			PointerValue *int

			// Unsafe Pointer - 1 / 0
			UnsafePointerValue unsafe.Pointer

			// unexported field - 1 / 0
			unexported int

			// Structs - top 6 / fields 8 / settable 6
			AnonStructValue struct{ Field int }
			StructValue     TestStruct
			TestStruct                              // exported embedded struct
			testStruct                              // unexported embedded struct
			PtrStruct       *struct{ Hello string } // initialize me
			NilStruct       *TestStruct             // leave me nil
		}

		t.Run("returns only settable fields", func(t *testing.T) {
			s := AllTypesStruct{
				ArrayValue:         [1]int{},
				SliceValue:         make([]int, 0),
				MapValue:           make(map[string]int),
				ChanValue:          make(chan int),
				FuncValue:          func() int { return 1 },
				PointerValue:       new(int),
				UnsafePointerValue: unsafe.Pointer(new(int)),
				AnonStructValue:    struct{ Field int }{Field: 1},
				PtrStruct:          &struct{ Hello string }{Hello: "World"},
			}

			fields, err := stronf.SettableFields(&s)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if testing.Verbose() {
				w := tabwriter.NewWriter(os.Stderr, 0, 0, 1, ' ', tabwriter.Debug)
				for i, f := range fields {
					fmt.Fprintf(w, "%d:\tname: %s\tkind: %s\ttype: %s\tvalue: %s\n", i+1, f.Name(), f.Kind(), f.Type(), f.Value())
				}
				w.Flush()
			}

			if len(fields) != 25 {
				t.Errorf("expected len 25, got len: %d", len(fields))
			}
		})
	}

	t.Run("unexported struct", func(t *testing.T) {
		type unexported struct {
			Int int
		}

		u := unexported{Int: 22}

		fields, err := stronf.SettableFields(&u)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d", len(fields))
		}
	})

	t.Run("anonymous struct", func(t *testing.T) {
		a := struct{ Int int }{Int: 55}

		fields, err := stronf.SettableFields(&a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d", len(fields))
		}
	})
}
