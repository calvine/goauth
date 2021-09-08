package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeTokenNotFound expired or token was found with the given id.
const ErrCodeTokenNotFound = "TokenNotFound"

// NewTokenNotFoundError creates a new specific error
func NewTokenNotFoundError(id string, includeStack bool) errors.RichError {
	msg := "expired or token was found with the given id."
	err := errors.NewRichError(ErrCodeTokenNotFound, msg).AddMetaData("id", id)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsTokenNotFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeTokenNotFound
}
