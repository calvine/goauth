package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidAppCreation this app had bad or missing required fields
const ErrCodeInvalidAppCreation = "InvalidAppCreation"

// NewInvalidAppCreationError creates a new specific error
func NewInvalidAppCreationError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "this app had bad or missing required fields"
	err := errors.NewRichError(ErrCodeInvalidAppCreation, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidAppCreationError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidAppCreation
}
