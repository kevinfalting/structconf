package confhandler

import (
	"context"
	"fmt"

	"github.com/kevinfalting/structconf/stronf"
)

// Required is a handler that will require a value to be set. Zero values are
// allowed when returned from handlers. It's typically best used as the very
// last handler.
type Required struct{}

func (h Required) Handle(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
	if proposedValue != nil {
		return proposedValue, nil
	}

	_, required := field.LookupTag("conf", "required")
	if required && field.IsZero() {
		return nil, fmt.Errorf("structconf: required field %q is not set", field.Name())
	}

	return nil, nil
}
