package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoContactFound no contact found for given query
const ErrCodeNoContactFound = "NoContactFound"

// NewNoContactFoundError creates a new specific error
func NewNoContactFoundError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "no contact found for given query"
	err := errors.NewRichError(ErrCodeNoContactFound, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoContactFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoContactFound
}
