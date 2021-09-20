package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidScopeCreation this scope had bad or missing required fields
const ErrCodeInvalidScopeCreation = "InvalidScopeCreation"

// NewInvalidScopeCreationError creates a new specific error
func NewInvalidScopeCreationError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "this scope had bad or missing required fields"
	err := errors.NewRichError(ErrCodeInvalidScopeCreation, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidScopeCreationError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidScopeCreation
}
