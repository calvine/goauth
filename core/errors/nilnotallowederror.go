package errors

import (
	"github.com/calvine/goauth/core/errors/codes"
)

func NewNilNotAllowedError() RichError {
	err := NewRichError(codes.ErrCodeNilNotAllowed, "a nil value was encountered, but not allowed")
	return err
}
