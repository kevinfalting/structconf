package stronf_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/kevinfalting/structconf/stronf"
)

func TestField_LookupTag(t *testing.T) {
	t.Run("key:\"tag:val,tag1:val1\"", func(t *testing.T) {
		A := struct {
			Int int `key:"tag:val,tag1:val1"`
		}{}

		fields, err := stronf.SettableFields(&A)
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

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key")
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

		fields, err := stronf.SettableFields(&A)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(fields) != 1 {
			t.Fatalf("expected fields len 1, got %d", len(fields))
		}

		field := fields[0]
		v, ok := field.LookupTag("key")
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

		fields, err := stronf.SettableFields(&A)
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

func ExampleField_LookupTag() {
	type Config struct {
		DatabaseURL string `conf:"env:DATABASE_URL,required,key:val" custom:"flag:db-url" emptyKey:""`
	}

	var config Config
	fields, err := stronf.SettableFields(&config)
	if err != nil {
		log.Println(err)
		return
	}

	if len(fields) != 1 {
		log.Println("expected 1 field")
		return
	}

	field := fields[0]
	if _, ok := field.LookupTag("no-exist", "hello"); ok {
		log.Println("should not recieve a value for a key that doesn't exist")
	}

	if val, ok := field.LookupTag("emptyKey"); ok {
		fmt.Printf("val: %q\n", val)
	}

	if val, ok := field.LookupTag("custom", "flag"); ok {
		fmt.Printf("val: %q\n", val)
	}

	if val, ok := field.LookupTag("conf", "env"); ok {
		fmt.Printf("val: %q\n", val)
	}

	if val, ok := field.LookupTag("conf", "required"); ok {
		fmt.Printf("val: %q\n", val)
	}
	// Output:
	// val: ""
	// val: "db-url"
	// val: "DATABASE_URL"
	// val: ""
}

func TestField_Parse(t *testing.T) {
	type Config struct {
		String string
	}

	t.Run("nil handler", func(t *testing.T) {
		var field stronf.Field
		err := field.Parse(context.Background(), nil)
		if err == nil {
			t.Error("expected error, got none")
		}
	})

	t.Run("should return an error if handler returns error and not update the field", func(t *testing.T) {
		var handler stronf.HandleFunc = func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
			return nil, errors.New("an error")
		}

		expect := "i'm a default value"

		config := Config{
			String: expect,
		}

		fields, err := stronf.SettableFields(&config)
		if err != nil {
			t.Fatal("failed to get SettableFields:", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		if err := fields[0].Parse(context.Background(), handler); err == nil {
			t.Error("expected error, got none")
		}

		got, ok := fields[0].Value().(string)
		if !ok {
			t.Fatalf("expected field value to be a string, got %T", fields[0].Value())
		}

		if got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})

	t.Run("should not override struct value if handler returns nil", func(t *testing.T) {
		var handler stronf.HandleFunc = func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
			// This additionally ensures that Field_Parse is passing in <nil> as the
			// initial proposedValue
			return proposedValue, nil
		}

		expect := "i'm a default value"

		config := Config{
			String: expect,
		}

		fields, err := stronf.SettableFields(&config)
		if err != nil {
			t.Fatal("failed to get SettableFields:", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		if err := fields[0].Parse(context.Background(), handler); err != nil {
			t.Fatal("failed to Parse:", err)
		}

		got, ok := fields[0].Value().(string)
		if !ok {
			t.Fatalf("expected field value to be a string, got %T", fields[0].Value())
		}

		if got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})

	t.Run("should set value returned by handler", func(t *testing.T) {
		expect := "i'm a value returned by the handler"
		var handler stronf.HandleFunc = func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
			return expect, nil
		}

		config := Config{
			String: "i'm a default value",
		}

		fields, err := stronf.SettableFields(&config)
		if err != nil {
			t.Fatal("failed to get SettableFields:", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		if err := fields[0].Parse(context.Background(), handler); err != nil {
			t.Fatal("failed to Parse:", err)
		}

		got, ok := fields[0].Value().(string)
		if !ok {
			t.Fatalf("expected field value to be a string, got %T", fields[0].Value())
		}

		if got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})

	t.Run("should error if value returned by handler is wrong type", func(t *testing.T) {
		var handler stronf.HandleFunc = func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
			return 55, nil
		}

		expect := "i'm a default value"

		config := Config{
			String: expect,
		}

		fields, err := stronf.SettableFields(&config)
		if err != nil {
			t.Fatal("failed to get SettableFields:", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		if err := fields[0].Parse(context.Background(), handler); err == nil {
			t.Error("expected error got none")
		}

		got, ok := fields[0].Value().(string)
		if !ok {
			t.Fatalf("expected field value to be a string, got %T", fields[0].Value())
		}

		if got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})

	t.Run("can set using unmarshalerFunc", func(t *testing.T) {
		type TimeConfig struct {
			Time time.Time
		}

		expectTimeString := "2023-11-05T15:04:05Z"
		expect, err := time.Parse(time.RFC3339, expectTimeString)
		if err != nil {
			t.Fatal("failed to Parse:", err)
		}

		var handler stronf.HandleFunc = func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
			return []byte(expectTimeString), nil
		}

		var config TimeConfig
		fields, err := stronf.SettableFields(&config)
		if err != nil {
			t.Fatal("failed to SettableFields:", err)
		}

		if len(fields) != 1 {
			t.Fatalf("expected 1 field, got %d", len(fields))
		}

		if err := fields[0].Parse(context.Background(), handler); err != nil {
			t.Error("failed to Parse:", err)
		}

		if config.Time.Compare(expect) != 0 {
			t.Errorf("expected %s, got %s", expect, config.Time)
		}
	})
}
