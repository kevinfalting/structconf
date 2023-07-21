package structconf_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf"
)

func TestDefault(t *testing.T) {
	t.Run("empty default tag", func(t *testing.T) {
		A := struct {
			Int int `conf:"default"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := structconf.HandlerFunc(
			func(ctx context.Context, f structconf.Field) (any, error) {
				return nil, nil
			},
		)

		h := structconf.WrapMiddleware([]structconf.Handler{handler}, structconf.Default())

		result, err := h.Handle(context.Background(), fields[0])
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected result to be nil, got %+v", result)
		}
	})

	t.Run("default tag applies when zero", func(t *testing.T) {
		A := struct {
			Int int `conf:"default:5"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := structconf.HandlerFunc(
			func(ctx context.Context, f structconf.Field) (any, error) {
				return nil, nil
			},
		)

		h := structconf.WrapMiddleware([]structconf.Handler{handler}, structconf.Default())

		result, err := h.Handle(context.Background(), fields[0])
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(5, result) {
			t.Errorf("expected 5, got %+v", result)
		}
	})

	t.Run("handler returns error, return error immediately", func(t *testing.T) {
		A := struct {
			Int int `conf:"default:5"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := structconf.HandlerFunc(
			func(ctx context.Context, f structconf.Field) (any, error) {
				return nil, errors.New("i'm an error")
			},
		)

		h := structconf.WrapMiddleware([]structconf.Handler{handler}, structconf.Default())

		result, err := h.Handle(context.Background(), fields[0])
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}
		if result != nil {
			t.Errorf("expected result to be nil, got %+v", result)
		}
	})

	t.Run("default tag does not apply when field is non-zero", func(t *testing.T) {
		A := struct {
			Int int `conf:"default:5"`
		}{
			Int: 88,
		}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := structconf.HandlerFunc(
			func(ctx context.Context, f structconf.Field) (any, error) {
				return nil, nil
			},
		)

		h := structconf.WrapMiddleware([]structconf.Handler{handler}, structconf.Default())

		result, err := h.Handle(context.Background(), fields[0])
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(88, A.Int) {
			t.Errorf("expected field to be 88, got %+v", result)
		}
	})

	t.Run("default does not apply when handler returns value", func(t *testing.T) {
		A := struct {
			Int int `conf:"default:5"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		handler := structconf.HandlerFunc(
			func(ctx context.Context, f structconf.Field) (any, error) {
				return 88, nil
			},
		)

		h := structconf.WrapMiddleware([]structconf.Handler{handler}, structconf.Default())

		result, err := h.Handle(context.Background(), fields[0])
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(result, 88) {
			t.Errorf("expected field to be 88, got %+v (%T)", result, result)
		}
	})
}
