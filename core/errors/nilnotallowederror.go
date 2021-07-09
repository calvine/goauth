package errors

import (
	"github.com/calvine/goauth/core/errors/codes"
)

func NewNilNotAllowedError(includeStack bool) RichError {
	msg := "a nil value was encountered, but not allowed"
	err := NewRichError(codes.ErrCodeNilNotAllowed, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
