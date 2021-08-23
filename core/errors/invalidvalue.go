package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidValue an invalid value was found
const ErrCodeInvalidValue = "InvalidValue"

// NewInvalidValueError creates a new specific error
func NewInvalidValueError(value interface{}, includeStack bool) errors.RichError {
	msg := "an invalid value was found"
	err := errors.NewRichError(ErrCodeInvalidValue, msg).AddMetaData("value", value)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidValueError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidValue
}
