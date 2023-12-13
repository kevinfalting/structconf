package confhandler

import (
	"context"
	"fmt"

	"github.com/kevinfalting/structconf/stronf"
)

// Default is a handler which will set the default value for the field as long
// as no value is being proposed and the field is the zero value for it's type.
// It's typically best just before the Required handler.
type Default struct{}

func (h Default) Handle(ctx context.Context, field stronf.Field, interimValue any) (any, error) {
	if interimValue != nil {
		return interimValue, nil
	}

	if !field.IsZero() {
		return interimValue, nil
	}

	defaultVal, ok := field.LookupTag("conf", "default")
	if !ok {
		return interimValue, nil
	}

	if len(defaultVal) == 0 {
		return nil, fmt.Errorf("structconf: empty default value for field %q", field.Name())
	}

	return defaultVal, nil
}
