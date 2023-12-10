package stronf_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/kevinfalting/structconf/stronf"
)

func TestWrapMiddleware(t *testing.T) {
	t.Run("nil wrapping", func(t *testing.T) {
		if h := stronf.WrapMiddleware(nil, nil); h != nil {
			t.Error("expected nil handler")
		}
	})

	t.Run("middleware is called in order", func(t *testing.T) {
		var mw1Time time.Time
		mw1 := func(h stronf.Handler) stronf.Handler {
			return stronf.HandlerFunc(func(ctx context.Context, f stronf.Field, a any) (any, error) {
				mw1Time = time.Now()
				h.Handle(ctx, f, a)
				return nil, nil
			})
		}

		var mw2Time time.Time
		mw2 := func(h stronf.Handler) stronf.Handler {
			return stronf.HandlerFunc(func(ctx context.Context, f stronf.Field, a any) (any, error) {
				mw2Time = time.Now()
				h.Handle(ctx, f, a)
				return nil, nil
			})
		}

		doNothingHandler := stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) { return nil, nil })

		h := stronf.WrapMiddleware(doNothingHandler, mw1, nil, mw2)
		if h == nil {
			t.Error("expected non-nil handler")
		}

		result, err := h.Handle(context.Background(), stronf.Field{}, nil)
		if err != nil {
			t.Error("failed to Handle:", err)
		}

		if result != nil {
			t.Error("expected result to be nil")
		}

		if !mw1Time.Before(mw2Time) {
			t.Errorf("expected %q to be before %q", mw1Time, mw2Time)
		}
	})

	t.Run("just the handler", func(t *testing.T) {
		doNothingHandler := stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, interimValue any) (any, error) { return nil, nil })

		h := stronf.WrapMiddleware(doNothingHandler, nil)

		ptr1 := reflect.ValueOf(doNothingHandler).Pointer()
		ptr2 := reflect.ValueOf(h).Pointer()

		if ptr1 != ptr2 {
			t.Error("did not return the same function")
		}
	})

	t.Run("just middleware", func(t *testing.T) {
		mw1 := func(h stronf.Handler) stronf.Handler {
			return stronf.HandlerFunc(func(ctx context.Context, f stronf.Field, a any) (any, error) {
				h.Handle(ctx, f, a)
				return nil, nil
			})
		}

		h := stronf.WrapMiddleware(nil, mw1)
		if h == nil {
			t.Error("wrapped handler should not be nil")
		}

		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic, but it didn't")
			}
		}()

		h.Handle(context.Background(), stronf.Field{}, nil)
	})
}
