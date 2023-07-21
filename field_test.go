package structconf_test

import (
	"testing"

	"github.com/kevinfalting/structconf"
)

func TestLookupTag(t *testing.T) {
	t.Run("key:\"tag:val,tag1:val1\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag:val,tag1:val1"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "tag")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val" {
			t.Errorf("expected v = val, got %q", v)
		}
	})

	t.Run("key:\"\"", func(t *testing.T) {
		A := struct {
			Int int `key:""`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "")
		if !ok {
			t.Error("expected key to exist")
		}

		if v != "" {
			t.Errorf("expected v = \"\", got %q", v)
		}
	})

	t.Run("key:\"tag\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "")
		if !ok {
			t.Error("expected key to exist")
		}

		if v != "tag" {
			t.Errorf("expected v = tag, got %q", v)
		}
	})

	t.Run("key:\"tag:val\" key1:\"tag1:val1\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag:val" key1:"tag1:val1"`
		}{}

		fields, err := structconf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key", "tag")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val" {
			t.Errorf("expected v = val, got %q", v)
		}

		v, ok = field.LookupTag("key1", "tag1")
		if !ok {
			t.Error("expected key tag to exist")
		}

		if v != "val1" {
			t.Errorf("expected v = val, got %q", v)
		}
	})
}
