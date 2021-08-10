package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewInvalidValueError creates a new specific error
func NewInvalidValueError(value interface{}, includeStack bool) RichError {
	msg := "an invalid value was found"
	err := NewRichError(codes.ErrCodeInvalidValue, msg).AddMetaData("value", value)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
