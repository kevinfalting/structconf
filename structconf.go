package structconf

import (
	"context"
	"flag"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

// Parse will set any settable fields in the provided struct based on the
// results of the default handlers. By default, it checks for environment
// variables, default, and required tags. Flags are optionally enabled.
func Parse(ctx context.Context, cfg any, optionFuncs ...optionFunc) error {
	fields, err := stronf.SettableFields(cfg)
	if err != nil {
		return err
	}

	var opt option
	for _, optionFunc := range optionFuncs {
		optionFunc(&opt)
	}

	handlers := []stronf.HandleFunc{
		confhandler.EnvironmentVariable{}.Handle,
	}

	if opt.flagSet != nil {
		flagHandler := confhandler.NewFlag(opt.flagSet)
		if err := flagHandler.DefineFlags(fields); err != nil {
			return err
		}

		handlers = append(handlers, flagHandler.Handle)
	}

	handlers = append(handlers,
		confhandler.Default{}.Handle,
		confhandler.Required{}.Handle,
	)

	handler := stronf.CombineHandlers(handlers...)

	for _, field := range fields {
		if err := field.Parse(ctx, handler); err != nil {
			return err
		}
	}

	return nil
}

type option struct {
	flagSet *flag.FlagSet
}

type optionFunc func(opt *option)

// WithFlagSet will signal to use the [confhandler.Flag] and optionally pass it
// a [flag.FlagSet]. Passing a nil flagset will create and use one similar to
// the stdlib flagset.
func WithFlagSet(fset *flag.FlagSet) optionFunc {
	return func(opt *option) {
		opt.flagSet = fset
	}
}
