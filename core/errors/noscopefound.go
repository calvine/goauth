package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoScopeFound no scope found for given query
const ErrCodeNoScopeFound = "NoScopeFound"

// NewNoScopeFoundError creates a new specific error
func NewNoScopeFoundError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "no scope found for given query"
	err := errors.NewRichError(ErrCodeNoScopeFound, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoScopeFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoScopeFound
}
