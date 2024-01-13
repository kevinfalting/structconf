package structconf

import (
	"context"
	"flag"

	"github.com/kevinfalting/structconf/confhandler"
	"github.com/kevinfalting/structconf/stronf"
)

// Parse will set any settable fields in the provided struct based on the
// results of the default handlers.
func Parse(ctx context.Context, cfg any, optionFuncs ...optionFunc) error {
	fields, err := stronf.SettableFields(cfg)
	if err != nil {
		return err
	}

	var opt option
	for _, optionFunc := range optionFuncs {
		optionFunc(&opt)
	}

	flagHandler := confhandler.NewFlag(opt.flagSet)
	if err := flagHandler.DefineFlags(fields); err != nil {
		return err
	}

	handler := stronf.CombineHandlers(
		confhandler.EnvironmentVariable{}.Handle,
		flagHandler.Handle,
		confhandler.Default{}.Handle,
		confhandler.Required{}.Handle,
	)

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

// WithFlagSet is a functional option for passing a [flag.FlagSet] to the flag
// handler.
func WithFlagSet(fset *flag.FlagSet) optionFunc {
	return func(opt *option) {
		opt.flagSet = fset
	}
}
