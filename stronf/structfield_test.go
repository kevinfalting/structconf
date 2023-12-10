package stronf_test

import (
	"encoding"
	"math/big"
	"net"
	"net/netip"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/kevinfalting/structconf/stronf"
)

func TestSettableFields(t *testing.T) {
	t.Run("supported field types", func(t *testing.T) {
		type StdlibTypes struct {
			// These types are supported because they implement either
			// encoding.TextUnmarshaler or (inclusive) encoding.UnmarshalBinary.
			// This isn't exhaustive, just making explicit the implicit.
			BigFloat big.Float
			BigInt   big.Int
			BigRat   big.Rat

			NetIP net.IP

			NetIPAddr     netip.Addr
			NetIPAddrPort netip.AddrPort
			NetIPPrefix   netip.Prefix

			Regex regexp.Regexp

			TimeDuration time.Duration
			Time         time.Time

			URL url.URL
		}

		type BuiltinTypes struct {
			Bool bool

			Complex64  complex64
			Complex128 complex128

			Float32 float32
			Float64 float64

			Int   int
			Int8  int8
			Int16 int16
			Int32 int32
			Int64 int64

			String string

			Uint   uint
			Uint8  uint8
			Uint16 uint16
			Uint32 uint32
			Uint64 uint64

			Uintptr uintptr

			NestedStruct StdlibTypes
		}

		x := BuiltinTypes{}
		fields, err := stronf.SettableFields(&x)
		if err != nil {
			t.Fatal("failed to SettableFields:", err)
		}

		expectNumberOfFields := 27
		if len(fields) != expectNumberOfFields {
			for _, field := range fields {
				t.Errorf("field name: %q, value: %q", field.Name(), field.Value())
			}
			t.Errorf("expected %d, got %d", expectNumberOfFields, len(fields))
		}
	})

	t.Run("unsupported field types", func(t *testing.T) {
		type StdlibTypes struct {
			// Just being specific about not supporting an interface field, even if
			// that field is the interface type being check for.
			BinaryUnmarshaler encoding.BinaryUnmarshaler
			TextUnmarshaler   encoding.TextUnmarshaler

			TimePtr *time.Time
		}

		type BuiltinTypes struct {
			Array     [0]string
			Channel   chan string
			Function  func()
			Interface interface{}
			Map       map[string]string
			Slice     []string
			Pointer   *struct{}

			NestedStruct StdlibTypes
		}

		x := BuiltinTypes{}
		fields, err := stronf.SettableFields(&x)
		if err != nil {
			t.Fatal("failed to SettableFields:", err)
		}

		if len(fields) != 0 {
			for _, field := range fields {
				t.Errorf("expected field %q, type %q to be unsupported", field.Name(), field.Type())
			}
			t.Errorf("expected 0, got %d", len(fields))
		}
	})

	t.Run("unmarshal text func", func(t *testing.T) {
		type X struct {
			Time time.Time
		}

		fields, err := stronf.SettableFields(&X{})
		if err != nil {
			t.Fatal("failed to SettableFields", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		testTime := "2023-11-05T15:04:05Z"

		expectTime, err := time.Parse(time.RFC3339, testTime)
		if err != nil {
			t.Fatalf("failed to parse %s: %s", testTime, err)
		}

		if err := fields[0].Set([]byte(testTime)); err != nil {
			t.Fatal("failed to Set", err)
		}

		gotTime, ok := fields[0].Value().(time.Time)
		if !ok {
			t.Fatal("failed to cast to time")
		}

		if expectTime.Compare(gotTime) != 0 {
			t.Errorf("expected %s, got %s", expectTime, gotTime)
		}
	})

	t.Run("handle time.Duration", func(t *testing.T) {
		type X struct {
			Duration time.Duration
		}

		fields, err := stronf.SettableFields(&X{})
		if err != nil {
			t.Fatal("failed to SettableFields", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		testDuration := "3h55m2s"
		expectDuration, err := time.ParseDuration(testDuration)
		if err != nil {
			t.Fatal("failed to parse duration", err)
		}

		if err := fields[0].Set(testDuration); err != nil {
			t.Fatal("failed to Set duration", err)
		}

		gotDuration, ok := fields[0].Value().(time.Duration)
		if !ok {
			t.Fatal("failed to cast to duration")
		}

		if gotDuration != expectDuration {
			t.Errorf("expected %s, got %s", expectDuration, gotDuration)
		}
	})
}
