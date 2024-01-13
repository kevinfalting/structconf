package confhandler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestRequired(t *testing.T) {
	type A struct {
		Int   int `conf:"required"`
		NoTag int
	}

	tests := map[string]struct {
		input         A
		proposedValue any
		expect        A
	}{
		"field value set should be okay": {
			input: A{
				Int: 2,
			},
			expect: A{
				Int: 2,
			},
		},
		"proposedValue should overrid all other values": {
			input: A{
				Int:   2,
				NoTag: 22,
			},
			proposedValue: 55,
			expect: A{
				Int:   55,
				NoTag: 55,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			requiredHandler := confhandler.Required{}
			proposedValueRequiredHandler := func(testProposedValue any) stronf.HandleFunc {
				return func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
					return requiredHandler.Handle(ctx, field, testProposedValue)
				}
			}(test.proposedValue)

			fields, err := stronf.SettableFields(&test.input)
			if err != nil {
				t.Fatal("failed to SettableFields:", err)
			}

			for _, field := range fields {
				if err := field.Parse(context.Background(), proposedValueRequiredHandler); err != nil {
					t.Error("expected no error, got:", err)
				}
			}

			if !reflect.DeepEqual(test.expect, test.input) {
				t.Errorf("\nexpected:\n%+v\ngot:\n%+v", test.expect, test.input)
			}
		})
	}

	t.Run("required and zero value with no proposedValue", func(t *testing.T) {
		type B struct {
			Int int `conf:"required"`
		}

		var b B
		fields, err := stronf.SettableFields(&b)
		if err != nil {
			t.Fatal("failed to SettableFields:", err)
		}

		err = fields[0].Parse(context.Background(), confhandler.Required{}.Handle)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
