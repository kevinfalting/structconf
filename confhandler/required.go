package confhandler

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kevinfalting/structconf/stronf"
)

// Required is a handler that will require a value to be set. Zero values are
// allowed when returned from handlers. It's typically best used as the very
// last handler.
type Required struct{}

func (h Required) Handle(ctx context.Context, field stronf.Field, proposedValue any) (any, error) {
	if proposedValue != nil && reflect.ValueOf(proposedValue).IsZero() {
		return proposedValue, nil
	}

	if !field.IsZero() {
		return proposedValue, nil
	}

	_, required := field.LookupTag("conf", "required")

	if required && proposedValue == nil {
		return nil, fmt.Errorf("structconf: required field %s is not set", field.Name())
	}

	return proposedValue, nil
}
