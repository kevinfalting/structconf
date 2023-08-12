package confhandler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestRequired(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			A := struct {
				Int int `conf:"required"`
			}{}

			fields, err := stronf.SettableFields(&A)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(fields) != 1 {
				t.Fatalf("expected 1 field, got %d", len(fields))
			}

			handler := stronf.HandlerFunc(
				func(ctx context.Context, f stronf.Field, _ any) (any, error) {
					return nil, nil
				},
			)

			h := stronf.WrapMiddleware([]stronf.Handler{handler}, confhandler.Required())

			result, err := h.Handle(context.Background(), fields[0], nil)
			if err == nil {
				t.Errorf("expected an error, got %v", err)
			}
			if result != nil {
				t.Errorf("expected result to be nil, got %+v", result)
			}
		})

		t.Run("0", func(t *testing.T) {
			A := struct {
				Int int `conf:"required"`
			}{}

			fields, err := stronf.SettableFields(&A)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(fields) != 1 {
				t.Fatalf("expected 1 field, got %d", len(fields))
			}

			handler := stronf.HandlerFunc(
				func(ctx context.Context, f stronf.Field, _ any) (any, error) {
					return 0, nil
				},
			)

			h := stronf.WrapMiddleware([]stronf.Handler{handler}, confhandler.Required())

			result, err := h.Handle(context.Background(), fields[0], nil)
			if err != nil {
				t.Errorf("expected no err, got %v", err)
			}
			if result == nil {
				t.Errorf("expected result to not be nil, got %+v", result)
			}
		})
	})

	t.Run("no tags, allow any with no error", func(t *testing.T) {
		t.Run("return nil", func(t *testing.T) {
			A := struct {
				Int int `conf:""`
			}{}

			fields, err := stronf.SettableFields(&A)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(fields) != 1 {
				t.Fatalf("expected 1 field, got %d", len(fields))
			}

			handler := stronf.HandlerFunc(
				func(ctx context.Context, f stronf.Field, _ any) (any, error) {
					return nil, nil
				},
			)

			h := stronf.WrapMiddleware([]stronf.Handler{handler}, confhandler.Required())

			result, err := h.Handle(context.Background(), fields[0], nil)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result != nil {
				t.Errorf("expected result to be nil, got %+v", result)
			}
		})

		t.Run("return 0", func(t *testing.T) {
			A := struct {
				Int int `conf:""`
			}{}

			fields, err := stronf.SettableFields(&A)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(fields) != 1 {
				t.Fatalf("expected 1 field, got %d", len(fields))
			}

			handler := stronf.HandlerFunc(
				func(ctx context.Context, f stronf.Field, _ any) (any, error) {
					return 0, nil
				},
			)

			h := stronf.WrapMiddleware([]stronf.Handler{handler}, confhandler.Required())

			result, err := h.Handle(context.Background(), fields[0], nil)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result == nil {
				t.Errorf("expected result to not be nil, got %+v", result)
			}
		})
	})

	t.Run("handler returns error, immediately return that error", func(t *testing.T) {
		A := struct {
			Int int `conf:""`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := stronf.HandlerFunc(
			func(ctx context.Context, f stronf.Field, _ any) (any, error) {
				return nil, errors.New("handler error")
			},
		)

		h := stronf.WrapMiddleware([]stronf.Handler{handler}, confhandler.Required())

		result, err := h.Handle(context.Background(), fields[0], nil)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected result to be nil, got %+v", result)
		}
	})
}

func TestIsZero(t *testing.T) {
	type TestStruct struct {
		A int
		B string
	}

	tests := []struct {
		input any
		want  bool
	}{
		{0, true},
		{1, false},
		{float32(0), true},
		{float32(1.1), false},
		{"", true},
		{"abc", false},
		{nil, true},
		{TestStruct{}, true},
		{TestStruct{A: 1, B: "abc"}, false},
		{[]int{}, true},
		{[]int{1, 2, 3}, false},
		{map[string]int{}, true},
		{map[string]int{"one": 1, "two": 2}, false},
		{[...]int{0, 0, 0}, true},
		{[...]int{1, 2, 3}, false},
		{func() {}, false},
	}

	for i, tt := range tests {
		got := confhandler.IsZero(tt.input)
		if got != tt.want {
			t.Errorf("Test %d: expected %v, got %v", i, tt.want, got)
		}
	}
}
