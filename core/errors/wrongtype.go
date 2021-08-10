package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewWrongTypeError creates a new specific error
func NewWrongTypeError(actual string, expected string, includeStack bool) RichError {
	msg := "unexpected type encountered"
	err := NewRichError(codes.ErrCodeWrongType, msg).AddMetaData("actual", actual).AddMetaData("expected", expected)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
