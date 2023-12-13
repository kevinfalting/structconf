package confhandler_test

import (
	"context"
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

			result, err := confhandler.Required{}.Handle(context.Background(), fields[0], nil)
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

			result, err := confhandler.Required{}.Handle(context.Background(), fields[0], 0)
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

			result, err := confhandler.Required{}.Handle(context.Background(), fields[0], nil)
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

			result, err := confhandler.Required{}.Handle(context.Background(), fields[0], 0)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result == nil {
				t.Errorf("expected result to not be nil, got %+v", result)
			}
		})
	})
}
