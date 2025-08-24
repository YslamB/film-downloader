package utils

import (
	"fmt"
	"log"
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

func LogAndReturnError(err error, context string) error {
	if err == nil {
		return nil
	}
	wrappedErr := WrapError(err, context)
	log.Printf("Error: %v", wrappedErr)
	return wrappedErr
}

func LogErrorf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	wrappedErr := WrapErrorf(err, format, args...)
	log.Printf("Error: %v", wrappedErr)
	return wrappedErr
}

func HandleCriticalError(err error, context string) {
	if err == nil {
		return
	}
	log.Fatalf("Critical error in %s: %v", context, err)
}

func HandleCriticalErrorf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	context := fmt.Sprintf(format, args...)
	log.Fatalf("Critical error in %s: %v", context, err)
}

func CheckError(err error, message string) {
	if err != nil {
		log.Fatalf("❌ %s: %v", message, err)
	}
}

type RetryableError struct {
	Err       error
	Retryable bool
}

func (re *RetryableError) Error() string {
	return re.Err.Error()
}

func (re *RetryableError) Unwrap() error {
	return re.Err
}

func NewRetryableError(err error, retryable bool) *RetryableError {
	return &RetryableError{
		Err:       err,
		Retryable: retryable,
	}
}

func IsRetryable(err error) bool {
	if retryableErr, ok := err.(*RetryableError); ok {
		return retryableErr.Retryable
	}
	return false
}
