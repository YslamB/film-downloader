package utils

import (
	"fmt"
)

func WrapError(err error, context string) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

func WrapErrorf(err error, format string, args ...interface{}) error {

	if err == nil {
		return nil
	}
	context := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", context, err)
}
