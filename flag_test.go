package structconf_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf"
)

func TestFlags(t *testing.T) {
	type A struct {
		Int int `conf:"flag:int"`
	}

	t.Run("non-zero flag value provided", func(t *testing.T) {
		var a A
		fields, err := structconf.SettableFields(&a)
		if err != nil {
			t.Fatalf("failed to SettableFields: %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		flagsHandler, err := structconf.NewFlag[A](nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = flagsHandler.Parse("-int", "5")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		result, err := flagsHandler.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(5, result) {
			t.Errorf("expected 5, got %+v", result)
		}
	})

	t.Run("zero flag value provided", func(t *testing.T) {
		var a A
		fields, err := structconf.SettableFields(&a)
		if err != nil {
			t.Fatalf("failed to SettableFields: %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		flagsHandler, err := structconf.NewFlag[A](nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = flagsHandler.Parse("-int", "0")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		result, err := flagsHandler.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(0, result) {
			t.Errorf("expected 0, got %+v", result)
		}
	})

	t.Run("no flag value provided", func(t *testing.T) {
		var a A
		fields, err := structconf.SettableFields(&a)
		if err != nil {
			t.Fatalf("failed to SettableFields: %v", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		flagsHandler, err := structconf.NewFlag[A](nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = flagsHandler.Parse("")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		result, err := flagsHandler.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(nil, result) {
			t.Errorf("expected nil, got %+v", result)
		}
	})
}
