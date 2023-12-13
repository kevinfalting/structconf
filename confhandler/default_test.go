package confhandler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestDefault(t *testing.T) {
	t.Run("empty default tag", func(t *testing.T) {
		A := struct {
			Int int `conf:"default"`
		}{}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		result, err := confhandler.Default{}.Handle(context.Background(), fields[0], nil)
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

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		result, err := confhandler.Default{}.Handle(context.Background(), fields[0], nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != "5" {
			t.Errorf("expected result to be 5, got %#v", result)
		}
	})

	t.Run("default tag does not apply when field is non-zero", func(t *testing.T) {
		A := struct {
			Int int `conf:"default:5"`
		}{
			Int: 88,
		}

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		result, err := confhandler.Default{}.Handle(context.Background(), fields[0], nil)
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

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		result, err := confhandler.Default{}.Handle(context.Background(), fields[0], 88)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !reflect.DeepEqual(result, 88) {
			t.Errorf("expected field to be 88, got %+v (%T)", result, result)
		}
	})
}
