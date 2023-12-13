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

func (h Required) Handle(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
	if interimValue != nil && reflect.ValueOf(interimValue).IsZero() {
		return interimValue, nil
	}

	if !field.IsZero() {
		return interimValue, nil
	}

	_, required := field.LookupTag("conf", "required")

	if required && interimValue == nil {
		return nil, fmt.Errorf("structconf: required field %s is not set", field.Name())
	}

	return interimValue, nil
}
