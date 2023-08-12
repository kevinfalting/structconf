package confhandler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestEnvironmentVariable(t *testing.T) {
	type A struct {
		Int int `conf:"env:INT"`
	}

	t.Run("env var set", func(t *testing.T) {
		var a A

		fields, err := stronf.SettableFields(&a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		t.Setenv("INT", "5")

		var ev confhandler.EnvironmentVariable

		result, err := ev.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result == nil {
			t.Errorf("expected not nil result, got %+v", result)
		}

		if !reflect.DeepEqual(5, result) {
			t.Errorf("expected 5, got %+v", result)
		}
	})

	t.Run("no env var set", func(t *testing.T) {
		var a A

		fields, err := stronf.SettableFields(&a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		var ev confhandler.EnvironmentVariable

		result, err := ev.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})

	t.Run("env var set wrong type - empty string", func(t *testing.T) {
		var a A

		fields, err := stronf.SettableFields(&a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		t.Setenv("INT", "")

		var ev confhandler.EnvironmentVariable

		result, err := ev.Handle(context.Background(), fields[0], nil)
		if err == nil {
			t.Errorf("expected error, got %v", err)
		}

		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})

	t.Run("field has no env tag", func(t *testing.T) {
		b := struct {
			Int int `conf:"asdf"`
		}{}

		fields, err := stronf.SettableFields(&b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		var ev confhandler.EnvironmentVariable

		result, err := ev.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != nil {
			t.Errorf("expected nil result, got %+v", result)
		}
	})
}
