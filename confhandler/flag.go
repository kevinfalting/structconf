package confhandler

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/kevinfalting/structconf/stronf"
)

// Flag maintains the state of the flag handler, providing the methods required
// to interact with it.
type Flag struct {
	fset *flag.FlagSet
}

// NewFlag returns an initialized [Flag] with the provided [flag.FlagSet]. If no
// [flag.FlagSet] is provided, one is initialized the same way the stdlib does.
func NewFlag(fset *flag.FlagSet) *Flag {
	if fset == nil {
		fset = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	flag := Flag{
		fset: fset,
	}

	return &flag
}

// Handle will lookup and set a field's flag value if one was set.
func (f *Flag) Handle(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
	if !f.fset.Parsed() {
		if err := f.Parse(os.Args[1:]); err != nil {
			return nil, err
		}
	}

	flagName, ok := field.LookupTag("conf", "flag")
	if !ok {
		return proposedValue, nil
	}

	stdFlag := f.fset.Lookup(flagName)
	if stdFlag == nil || stdFlag.Value == nil {
		return proposedValue, nil
	}

	fVal, ok := stdFlag.Value.(*flagVal)
	if !ok {
		return proposedValue, nil
	}

	if fVal.val != nil {
		return fVal.val, nil
	}

	if proposedValue != nil {
		return proposedValue, nil
	}

	return stdFlag.Value.String(), nil
}

// Parse passes the args to the underlying [flag.FlagSet]'s Parse method.
func (f *Flag) Parse(args []string) error {
	return f.fset.Parse(args)
}

// DefineFlags will define any flags on the [Flag]'s underlying [flag.FlagSet]
// based on the [stronf.Field]'s that are passed in. It looks for the "flag" tag
// first for the name, the looks up the "default" tag for the default value. If
// no "default" flag is provided, the [stronf.Field]'s value is used.
// Optionally, a "usage" flag can be provided to provide custom usage
// information.
func (f *Flag) DefineFlags(fields []stronf.Field) error {
	for _, field := range fields {
		if err := f.defineFlag(field); err != nil {
			return err
		}
	}

	return nil
}

func (f *Flag) defineFlag(field stronf.Field) error {
	flagName, ok := field.LookupTag("conf", "flag")
	if !ok {
		return nil
	}

	defaultVal, ok := field.LookupTag("conf", "default")
	if (!ok && field.Value() != nil) || !field.IsZero() {
		defaultVal = fmt.Sprintf("%v", field.Value())
	}

	usage, ok := field.LookupTag("conf", "usage")
	if !ok {
		usage = fmt.Sprintf("%s is a [`%T`]", flagName, field.Value())
	} else {
		usage = fmt.Sprintf("%s [`%T`]", usage, field.Value())
	}

	fVal := flagVal{
		field:      field,
		val:        nil,
		defaultVal: defaultVal,
	}

	f.fset.Var(&fVal, flagName, usage)

	return nil
}

var _ flag.Value = (*flagVal)(nil)

type flagVal struct {
	field      stronf.Field
	val        any
	defaultVal string
}

func (f *flagVal) Set(s string) error {
	val, err := stronf.Coerce(f.field, s)
	if err != nil {
		return err
	}

	f.val = val

	return nil
}

func (f *flagVal) String() string {
	if f.val == nil {
		return f.defaultVal
	}

	return fmt.Sprintf("%v", f.val)
}

func (f *flagVal) IsBoolFlag() bool {
	return f.field.Kind() == reflect.Bool
}
