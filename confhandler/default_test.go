package confhandler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestDefault(t *testing.T) {
	type A struct {
		Int   int `conf:"default:2"`
		NoTag int
	}

	tests := map[string]struct {
		input         A
		proposedValue any
		expect        A
	}{
		"nothing set should set to defaults": {
			expect: A{
				Int: 2,
			},
		},
		"proposedValue should be set no matter what else is set": {
			input: A{
				Int:   22,
				NoTag: 55,
			},
			proposedValue: 88,
			expect: A{
				Int:   88,
				NoTag: 88,
			},
		},
		"field value set takes precedence over default": {
			input: A{
				Int: 22,
			},
			expect: A{
				Int: 22,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var defaultHandler confhandler.Default
			proposedValueDefaultHandler := func(testProposedValue any) stronf.HandleFunc {
				return func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
					return defaultHandler.Handle(ctx, field, testProposedValue)
				}
			}(test.proposedValue)

			fields, err := stronf.SettableFields(&test.input)
			if err != nil {
				t.Fatal("failed to SettableFields:", err)
			}

			for _, field := range fields {
				if err := field.Parse(context.Background(), proposedValueDefaultHandler); err != nil {
					t.Error("expected no error, got:", err)
				}
			}

			if !reflect.DeepEqual(test.expect, test.input) {
				t.Errorf("\nexpected:\n%+v\ngot:\n%+v", test.expect, test.input)
			}
		})
	}

	t.Run("empty default value", func(t *testing.T) {
		type B struct {
			Int int `conf:"default"`
		}

		var b B
		fields, err := stronf.SettableFields(&b)
		if err != nil {
			t.Fatal("failed to SettableFields:", err)
		}

		err = fields[0].Parse(context.Background(), confhandler.Default{}.Handle)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
