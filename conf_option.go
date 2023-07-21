package structconf

import "flag"

type confOption struct {
	flagSet *flag.FlagSet
}

type confOptionFunc func(opt *confOption)

// WithFlagSet is a functional option for passing a [flag.FlagSet] to the
// default Conf's Flag handler.
func WithFlagSet(fset *flag.FlagSet) confOptionFunc {
	return func(opt *confOption) {
		opt.flagSet = fset
	}
}
