package stronf_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kevinfalting/structconf/stronf"
)

func TestCombineHandlers(t *testing.T) {
	// Test with no handlers (nil)
	t.Run("no handlers", func(t *testing.T) {
		h := stronf.CombineHandlers()
		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	// Test with one handler
	t.Run("one handler", func(t *testing.T) {
		h := stronf.CombineHandlers(stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
			return "one", nil
		}))
		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != "one" {
			t.Errorf("Expected 'one', got %v", result)
		}
	})

	// Test with two handlers
	t.Run("two handlers", func(t *testing.T) {
		h := stronf.CombineHandlers(
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "first", nil
			}),
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "second", nil
			}),
		)
		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != "second" {
			t.Errorf("Expected 'second', got %v", result)
		}
	})

	// Test error from handler returns immediately
	t.Run("error from handler", func(t *testing.T) {
		expectedError := errors.New("handler error")
		h := stronf.CombineHandlers(
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return nil, expectedError
			}),
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "should not be called", nil
			}),
		)
		_, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
	})

	// Test the result of one handler is passed into the interimValue of the next
	t.Run("pass interim value", func(t *testing.T) {
		h := stronf.CombineHandlers(
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "first", nil
			}),
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				if interimValue != "first" {
					t.Errorf("Expected 'first', got %v", interimValue)
				}
				return "second", nil
			}),
		)
		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != "second" {
			t.Errorf("Expected 'second', got %v", result)
		}
	})

	// Test the final interimValue is returned
	t.Run("final interim value", func(t *testing.T) {
		h := stronf.CombineHandlers(
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "first", nil
			}),
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				return "final", nil
			}),
		)
		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != "final" {
			t.Errorf("Expected 'final', got %v", result)
		}
	})

	// Test handlers are called in the order provided
	t.Run("handler order", func(t *testing.T) {
		var order []string
		h := stronf.CombineHandlers(
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				order = append(order, "first")
				return nil, nil
			}),
			stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
				order = append(order, "second")
				return nil, nil
			}),
		)
		_, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if order[0] != "first" || order[1] != "second" {
			t.Errorf("Handlers were not called in the correct order: got %v", order)
		}
	})
}
