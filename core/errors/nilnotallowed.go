package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewNilNotAllowedError creates a new specific error
func NewNilNotAllowedError(includeStack bool) RichError {
	msg := "a nil value was encountered, but not allowed"
	err := NewRichError(codes.ErrCodeNilNotAllowed, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
