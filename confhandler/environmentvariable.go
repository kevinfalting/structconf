package confhandler

import (
	"context"
	"os"

	"github.com/kevinfalting/structconf/stronf"
)

// EnvironmentVariable is a handler which will lookup in the environment for the
// 'env' key provided in the struct tag.
type EnvironmentVariable struct{}

var _ stronf.Handler = EnvironmentVariable{}

func (ev EnvironmentVariable) Handle(ctx context.Context, field stronf.Field, _ any) (any, error) {
	environmentVariable, ok := field.LookupTag("conf", "env")
	if !ok {
		return nil, nil
	}

	val, ok := os.LookupEnv(environmentVariable)
	if !ok {
		return nil, nil
	}

	return parseStringForKind(val, field.Kind())
}