package confhandler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

func TestFlags(t *testing.T) {
	type Config struct {
		Flag        int `conf:"flag:flag"`
		FlagDefault int `conf:"flag:flag-default,default:5"`
		NoFlag      int `conf:"default:2"` // with no flag defined, the flag handler will just use pass along the proposedValue, ignoring it otherwise
	}

	testCases := map[string]struct {
		input         Config
		args          []string
		proposedValue any
		expect        Config
	}{
		"no flags should result in defaults being set": {
			input:         Config{},
			args:          []string{},
			proposedValue: nil,
			expect: Config{
				FlagDefault: 5,
			},
		},
		"no flags with a proposedValue should use proposedValue": {
			input:         Config{},
			args:          []string{},
			proposedValue: 8,
			expect: Config{
				Flag:        8,
				FlagDefault: 8,
				NoFlag:      8,
			},
		},
		"flag provided should be used": {
			input:         Config{},
			args:          []string{"-flag=88", "-flag-default=55"},
			proposedValue: 22,
			expect: Config{
				Flag:        88,
				FlagDefault: 55,
				NoFlag:      22,
			},
		},
		"flag set to type's zero value should override any defaults": {
			input: Config{
				Flag:        88,
				FlagDefault: 55,
				NoFlag:      22,
			},
			args:          []string{"-flag=0", "-flag-default=0"},
			proposedValue: 258,
			expect: Config{
				NoFlag: 258,
			},
		},
		"no flag with default but the field value is set should leave the field value alone": {
			input: Config{
				FlagDefault: 88,
			},
			args:          nil,
			proposedValue: nil,
			expect: Config{
				FlagDefault: 88,
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			fields, err := stronf.SettableFields(&test.input)
			if err != nil {
				t.Fatalf("failed to SettableFields: %v", err)
			}

			flagsHandler := confhandler.NewFlag(nil)

			if err := flagsHandler.DefineFlags(fields); err != nil {
				t.Fatal("failed to DefineFlags:", err)
			}

			err = flagsHandler.Parse(test.args)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			proposedValueFlagHandler := func(testProposedValue any) stronf.Handler {
				return stronf.HandlerFunc(func(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
					return flagsHandler.Handle(ctx, field, testProposedValue)
				})
			}(test.proposedValue)

			for _, field := range fields {
				if err := field.Parse(context.Background(), proposedValueFlagHandler); err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}

			if !reflect.DeepEqual(test.expect, test.input) {
				t.Errorf("\nexpected:\n%+v\ngot:\n%+v", test.expect, test.input)
			}
		})
	}
}
