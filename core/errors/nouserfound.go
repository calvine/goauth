package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewNoUserFoundError creates a new specific error
func NewNoUserFoundError(fields map[string]interface{}, includeStack bool) RichError {
	msg := "no user found for given query"
	err := NewRichError(codes.ErrCodeNoUserFound, msg).AddMetaData("fields", fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
