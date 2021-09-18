package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNotImplemented functionality is not implemented
const ErrCodeNotImplemented = "NotImplemented"

// NewNotImplementedError creates a new specific error
func NewNotImplementedError(includeStack bool) errors.RichError {
	msg := "functionality is not implemented"
	err := errors.NewRichError(ErrCodeNotImplemented, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNotImplementedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNotImplemented
}
