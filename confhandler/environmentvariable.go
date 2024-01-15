package confhandler

import (
	"context"
	"os"

	"github.com/kevinfalting/structconf/stronf"
)

// EnvironmentVariable is a handler which will lookup in the environment for the
// 'env' key provided in the struct tag.
type EnvironmentVariable struct{}

// Handle is the [stronf.HandleFunc] implementation of the [EnvironmentVariable] handler.
func (ev EnvironmentVariable) Handle(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
	environmentVariable, ok := field.LookupTag("conf", "env")
	if !ok {
		return proposedValue, nil
	}

	val, ok := os.LookupEnv(environmentVariable)
	if !ok {
		return proposedValue, nil
	}

	return val, nil
}
