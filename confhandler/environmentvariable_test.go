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
		Int   int `conf:"env:INT"`
		NoTag int
	}

	testCases := map[string]struct {
		input         A
		proposedValue any
		envKey        string
		envVal        string
		expect        A
	}{
		"nothing set": {},
		"nothing set with proposedValue should set proposedValue": {
			proposedValue: 5,
			expect: A{
				Int:   5,
				NoTag: 5,
			},
		},
		"env set should set to env val": {
			envKey: "INT",
			envVal: "22",
			expect: A{
				Int: 22,
			},
		},
		"env set and proposedValue should set env": {
			envKey:        "INT",
			envVal:        "22",
			proposedValue: 5,
			expect: A{
				Int:   22,
				NoTag: 5,
			},
		},
		"env set and proposedValue with field value should ignore field value": {
			input: A{
				Int:   88,
				NoTag: 222,
			},
			envKey:        "INT",
			envVal:        "22",
			proposedValue: 5,
			expect: A{
				Int:   22,
				NoTag: 5,
			},
		},
		"only field values should be left alone": {
			input: A{
				Int:   88,
				NoTag: 222,
			},
			expect: A{
				Int:   88,
				NoTag: 222,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			fields, err := stronf.SettableFields(&test.input)
			if err != nil {
				t.Fatal("failed to SettableFields:", err)
			}

			var envHandler confhandler.EnvironmentVariable
			proposedValueEnvHandler := func(testProposedValue any) stronf.HandleFunc {
				return func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
					return envHandler.Handle(ctx, field, testProposedValue)
				}
			}(test.proposedValue)

			if len(test.envKey) != 0 {
				t.Setenv(test.envKey, test.envVal)
			}

			for _, field := range fields {
				if err := field.Parse(context.Background(), proposedValueEnvHandler); err != nil {
					t.Error("expected no error, got:", err)
				}
			}

			if !reflect.DeepEqual(test.expect, test.input) {
				t.Errorf("\nexpected:\n%+v\ngot:\n%+v", test.expect, test.input)
			}
		})
	}
}
