package structconf

import (
	"crypto/rsa"
	"flag"
)

type confOption struct {
	flagSet       *flag.FlagSet
	rsaPrivateKey *rsa.PrivateKey
	rsaLabel      []byte
}

type confOptionFunc func(opt *confOption)

// WithFlagSet is a functional option for passing a [flag.FlagSet] to the
// default Conf's Flag handler.
func WithFlagSet(fset *flag.FlagSet) confOptionFunc {
	return func(opt *confOption) {
		opt.flagSet = fset
	}
}

// WithRSAPrivateKey is a functional option for passing an [rsa.PrivateKey] to
// the default Conf's RSAHandler.
func WithRSAPrivateKey(priv *rsa.PrivateKey) confOptionFunc {
	return func(opt *confOption) {
		opt.rsaPrivateKey = priv
	}
}

// WithRSALabel sets the label to use with the RSA Private Key.
func WithRSALabel(label []byte) confOptionFunc {
	return func(opt *confOption) {
		opt.rsaLabel = label
	}
}
