package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidToken expired or token was found with the given id.
const ErrCodeInvalidToken = "InvalidToken"

// NewInvalidTokenError creates a new specific error
func NewInvalidTokenError(id string, includeStack bool) errors.RichError {
	msg := "expired or token was found with the given id."
	err := errors.NewRichError(ErrCodeInvalidToken, msg).AddMetaData("id", id)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidTokenError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidToken
}
