package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewInvalidTypeError creates a new specific error
func NewInvalidTypeError(typeEncountered string, includeStack bool) RichError {
	msg := "invalid type encountered"
	err := NewRichError(codes.ErrCodeInvalidType, msg).AddMetaData("typeEncountered", typeEncountered)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
