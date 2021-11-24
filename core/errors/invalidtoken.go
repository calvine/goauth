package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidToken expired or invalid token was found with the given id.
const ErrCodeInvalidToken = "InvalidToken"

// NewInvalidTokenError creates a new specific error
func NewInvalidTokenError(id string, includeStack bool) errors.RichError {
	msg := "expired or invalid token was found with the given id."
	err := errors.NewRichError(ErrCodeInvalidToken, msg).AddMetaData("id", id).WithTags([]string{"token"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidTokenError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidToken
}
