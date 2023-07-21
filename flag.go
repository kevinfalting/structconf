package structconf

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"
)

// Flag is the handler for parsing command line flags.
type Flag struct {
	flagSet     *flag.FlagSet
	parsedFlags map[string]func() any
}

// NewFlag returns an initialized Flag handler.
func NewFlag[T any](fset *flag.FlagSet) *Flag {
	if fset == nil {
		fset = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	a := new(T)

	fields, err := SettableFields(a)
	if err != nil {
		panic(err)
	}

	parsedFlags := make(map[string]func() any)

	fn := func(flagName string, kind reflect.Kind) func(string) error {
		return func(s string) error {
			val, err := parseStringForKind(s, kind)
			if err != nil {
				return err
			}

			parsedFlags[flagName] = func() any { return val }
			return nil
		}
	}

	for _, field := range fields {
		flagName, ok := field.LookupTag("conf", "flag")
		if !ok {
			continue
		}

		usage, _ := field.LookupTag("conf", "usage")

		switch field.Kind() {
		case reflect.Bool:
			var b bool
			db := defaultValFn[bool](field)
			fset.BoolVar(&b, flagName, db, usage)
			parsedFlags[flagName] = func() any { return valueFn[bool](b, db) }

		case reflect.Float64:
			var f float64
			df := defaultValFn[float64](field)
			fset.Float64Var(&f, flagName, df, usage)
			parsedFlags[flagName] = func() any { return valueFn[float64](f, df) }

		case reflect.Int64:
			var i int64
			di := defaultValFn[int64](field)
			fset.Int64Var(&i, flagName, di, usage)
			parsedFlags[flagName] = func() any { return valueFn[int64](i, di) }

		case reflect.Int:
			var i int
			di := defaultValFn[int](field)
			fset.IntVar(&i, flagName, defaultValFn[int](field), usage)
			parsedFlags[flagName] = func() any { return valueFn[int](i, di) }

		case reflect.String:
			var s string
			ds := defaultValFn[string](field)
			fset.StringVar(&s, flagName, ds, usage)
			parsedFlags[flagName] = func() any { return valueFn[string](s, ds) }

		case reflect.Uint64:
			var u uint64
			du := defaultValFn[uint64](field)
			fset.Uint64Var(&u, flagName, du, usage)
			parsedFlags[flagName] = func() any { return valueFn[uint64](u, du) }

		case reflect.Uint:
			var u uint
			du := defaultValFn[uint](field)
			fset.UintVar(&u, flagName, du, usage)
			parsedFlags[flagName] = func() any { return valueFn[uint](u, du) }

		default:
			fset.Func(flagName, usage, fn(flagName, field.Kind()))
		}
	}

	return &Flag{
		flagSet:     fset,
		parsedFlags: parsedFlags,
	}
}

var _ Handler = (*Flag)(nil)

func (f *Flag) Handle(ctx context.Context, field Field) (any, error) {
	flagName, ok := field.LookupTag("conf", "flag")
	if !ok {
		return nil, nil
	}

	if err := f.Parse(); err != nil {
		return nil, err
	}

	valueFn, ok := f.parsedFlags[flagName]
	if !ok {
		return nil, nil
	}

	return valueFn(), nil
}

// Parse will parse the flags associated with this Flag handler. If no args are
// provided, it will use os.Args. This is safe to call multiple times, first call
// wins.
func (f *Flag) Parse(args ...string) error {
	if f.flagSet.Parsed() {
		return nil
	}

	if len(args) == 0 {
		args = os.Args[1:]
	}

	return f.flagSet.Parse(args)
}

func defaultValFn[T any](f Field) T {
	defaultVal, ok := f.LookupTag("conf", "default")
	if !ok {
		return *new(T)
	}

	result, err := parseStringForKind(defaultVal, f.Kind())
	if err != nil {
		panic(err)
	}

	if _, ok := result.(T); !ok {
		panic(fmt.Errorf("value of type %T is not %T", result, *new(T)))
	}

	return result.(T)
}

func valueFn[T comparable](parsedVal, defaultVal T) any {
	if parsedVal == defaultVal {
		return nil
	}

	return parsedVal
}
