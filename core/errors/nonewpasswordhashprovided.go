package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoNewPasswordHashProvided no new password hash was provided
const ErrCodeNoNewPasswordHashProvided = "NoNewPasswordHashProvided"

// NewNoNewPasswordHashProvidedError creates a new specific error
func NewNoNewPasswordHashProvidedError(includeStack bool) errors.RichError {
	msg := "no new password hash was provided"
	err := errors.NewRichError(ErrCodeNoNewPasswordHashProvided, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoNewPasswordHashProvidedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoNewPasswordHashProvided
}
