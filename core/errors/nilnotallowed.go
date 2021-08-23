package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNilNotAllowed a nil value was encountered, but not allowed
const ErrCodeNilNotAllowed = "NilNotAllowed"

// NewNilNotAllowedError creates a new specific error
func NewNilNotAllowedError(includeStack bool) errors.RichError {
	msg := "a nil value was encountered, but not allowed"
	err := errors.NewRichError(ErrCodeNilNotAllowed, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNilNotAllowedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNilNotAllowed
}
