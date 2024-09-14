package types

import (
	"fmt"
)

func ValidateEpochIdentifierInterface(i any) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return ValidateEpochIdentifierString(v)
}

func ValidateEpochIdentifierString(s string) error {
	if s == "" {
		return fmt.Errorf("empty distribution epoch identifier: %+v", s)
	}

	return nil
}
