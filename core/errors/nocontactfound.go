package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewNoContactFoundError creates a new specific error
func NewNoContactFoundError(fields map[string]interface{}, includeStack bool) RichError {
	msg := "no contact found for given query"
	err := NewRichError(codes.ErrCodeNoContactFound, msg).AddMetaData("fields", fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
