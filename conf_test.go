package structconf_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
	"unsafe"

	"github.com/kevinfalting/structconf"
)

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
				_, err := structconf.SettableFields(tt.arg)
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
				_, err := structconf.SettableFields(tt.arg)
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

			fields, err := structconf.SettableFields(&s)
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

		fields, err := structconf.SettableFields(&u)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d", len(fields))
		}
	})

	t.Run("anonymous struct", func(t *testing.T) {
		a := struct{ Int int }{Int: 55}

		fields, err := structconf.SettableFields(&a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Errorf("expected 1 field, got %d", len(fields))
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("mw & handlers", func(t *testing.T) {
		type thing struct {
			Int int
		}

		t.Run("no middleware no handlers", func(t *testing.T) {
			a := thing{Int: 22}

			conf := structconf.Conf[thing]{}
			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}

			if a.Int != 22 {
				t.Errorf("expected 22, got %d", a.Int)
			}
		})

		t.Run("one handler, no middleware", func(t *testing.T) {
			a := thing{}

			conf := structconf.Conf[thing]{
				Handlers: []structconf.Handler{
					structconf.HandlerFunc(
						func(ctx context.Context, f structconf.Field, _ any) (any, error) {
							return int(5), nil
						},
					),
				},
			}

			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}

			if a.Int != 5 {
				t.Errorf("expected 5, got %d", a.Int)
			}
		})
		t.Run("one middleware, no handler", func(t *testing.T) {
			a := thing{}

			conf := structconf.Conf[thing]{
				Middleware: []structconf.Middleware{
					func(h structconf.Handler) structconf.Handler {
						return structconf.HandlerFunc(
							func(ctx context.Context, f structconf.Field, iv any) (any, error) {
								return h.Handle(ctx, f, iv)
							},
						)
					},
				},
			}

			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("default Conf", func(t *testing.T) {
		ctx := context.Background()

		t.Run("handler overwrites val in struct", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			t.Setenv("NAME", "22")

			a := A{
				Int: 88,
			}

			if err := conf.Parse(ctx, &a); err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if a.Int != 22 {
				t.Errorf("expected 22, got %d", a.Int)
			}
		})

		t.Run("handler returns nil, does not overwrite val in struct", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			a := A{
				Int: 88,
			}

			if err := conf.Parse(ctx, &a); err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if a.Int != 88 {
				t.Errorf("expected 88, got %d", a.Int)
			}
		})

		t.Run("required", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME,required"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			t.Run("set env val", func(t *testing.T) {
				t.Setenv("NAME", "5")

				var a A

				err := conf.Parse(ctx, &a)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 5 {
					t.Errorf("expected 5, got %d", a.Int)
				}
			})

			t.Run("no env val set but struct has a val set", func(t *testing.T) {
				a := A{
					Int: 8,
				}

				err := conf.Parse(ctx, &a)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 8 {
					t.Errorf("expected 8, got %d", a.Int)
				}
			})
		})

		t.Run("env & flag using the default tag", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:INT,flag:int,default:55"`
			}

			t.Run("environment variable set", func(t *testing.T) {
				var a A

				t.Setenv("INT", "22")

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse(nil); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("flag set", func(t *testing.T) {
				var a A

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse([]string{"-int", "22"}); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("flag & env set, expect flag to overwrite", func(t *testing.T) {
				var a A

				t.Setenv("INT", "456")

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse([]string{"-int", "22"}); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("no flag or env set, expect default to be set", func(t *testing.T) {
				var a A

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse(nil); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 55 {
					t.Errorf("expected 55, got %d", a.Int)
				}
			})
		})
	})
}
