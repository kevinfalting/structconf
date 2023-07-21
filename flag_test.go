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

	var a A
	fields, err := structconf.SettableFields(&a)
	if err != nil {
		t.Fatalf("failed to SettableFields: %v", err)
	}

	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}

	flagsHandler := structconf.NewFlag[A](nil)
	err = flagsHandler.Parse("-int", "5")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	result, err := flagsHandler.Handle(context.Background(), fields[0])
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(5, result) {
		t.Errorf("expected 5, got %+v", result)
	}
}
