package confhandler

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/kevinfalting/structconf/stronf"
)

// Flag is the handler for parsing command line flags.
type Flag struct {
	flagSet     *flag.FlagSet
	parsedFlags map[string]func() any
}

// NewFlag returns an initialized Flag handler.
func NewFlag[T any](fset *flag.FlagSet) (*Flag, error) {
	if fset == nil {
		fset = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	a := new(T)

	fields, err := stronf.SettableFields(a)
	if err != nil {
		return nil, err
	}

	parsedFlags := make(map[string]func() any)

	for _, field := range fields {
		flagName, ok := field.LookupTag("conf", "flag")
		if !ok {
			continue
		}

		usage, _ := field.LookupTag("conf", "usage")

		switch field.Kind() {
		case reflect.Bool:
			var b bool
			db, err := defaultValFn[bool](field)
			if err != nil {
				return nil, err
			}
			fset.BoolVar(&b, flagName, db, usage)
			parsedFlags[flagName] = func() any { return b }

		case reflect.Float64:
			var f float64
			df, err := defaultValFn[float64](field)
			if err != nil {
				return nil, err
			}
			fset.Float64Var(&f, flagName, df, usage)
			parsedFlags[flagName] = func() any { return f }

		case reflect.Int64:
			var i int64
			di, err := defaultValFn[int64](field)
			if err != nil {
				return nil, err
			}
			fset.Int64Var(&i, flagName, di, usage)
			parsedFlags[flagName] = func() any { return i }

		case reflect.Int:
			var i int
			di, err := defaultValFn[int](field)
			if err != nil {
				return nil, err
			}
			fset.IntVar(&i, flagName, di, usage)
			parsedFlags[flagName] = func() any { return i }

		case reflect.String:
			var s string
			ds, err := defaultValFn[string](field)
			if err != nil {
				return nil, err
			}
			fset.StringVar(&s, flagName, ds, usage)
			parsedFlags[flagName] = func() any { return s }

		case reflect.Uint64:
			var u uint64
			du, err := defaultValFn[uint64](field)
			if err != nil {
				return nil, err
			}
			fset.Uint64Var(&u, flagName, du, usage)
			parsedFlags[flagName] = func() any { return u }

		case reflect.Uint:
			var u uint
			du, err := defaultValFn[uint](field)
			if err != nil {
				return nil, err
			}
			fset.UintVar(&u, flagName, du, usage)
			parsedFlags[flagName] = func() any { return u }

		default:
			parseFlagValueFn := func(flagName string, field stronf.Field) func(string) error {
				return func(s string) error {
					parsedFlags[flagName] = func() any { return s }
					return nil
				}
			}

			fset.Func(flagName, usage, parseFlagValueFn(flagName, field))
		}
	}

	f := Flag{
		flagSet:     fset,
		parsedFlags: parsedFlags,
	}

	return &f, nil
}

var _ stronf.Handler = (*Flag)(nil)

func (f *Flag) Handle(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
	flagName, ok := field.LookupTag("conf", "flag")
	if !ok {
		return interimValue, nil
	}

	if err := f.Parse(); err != nil {
		return nil, err
	}

	valueFn, ok := f.parsedFlags[flagName]
	if !ok {
		return interimValue, nil
	}

	return valueFn(), nil
}

// Parse will parse the flags associated with this Flag handler. If no args are
// provided, it will use os.Args. This is safe to call multiple times, first call
// wins.
func (f *Flag) Parse(args ...string) error {
	if !f.flagSet.Parsed() {
		if len(args) == 0 {
			args = os.Args[1:]
		}

		if err := f.flagSet.Parse(args); err != nil {
			return err
		}
	}

	// TODO: This is called on every call to Handle, which isn't necessary. But,
	// is it a problem? Probably not. I'll optimize this when it's a problem.
	commandLineProvidedFlags := make(map[string]bool)

	f.flagSet.Visit(func(f *flag.Flag) {
		commandLineProvidedFlags[f.Name] = true
	})

	for flagName := range f.parsedFlags {
		if _, exists := commandLineProvidedFlags[flagName]; !exists {
			delete(f.parsedFlags, flagName)
		}
	}

	return nil
}

func defaultValFn[T any](field stronf.Field) (T, error) {
	defaultVal, ok := field.LookupTag("conf", "default")
	if !ok {
		return *new(T), nil
	}

	result, err := stronf.Coerce(field, defaultVal)
	if err != nil {
		return *new(T), nil
	}

	if _, ok := result.(T); !ok {
		return *new(T), fmt.Errorf("structconf: value of type %T is not %T", result, *new(T))
	}

	return result.(T), nil
}
