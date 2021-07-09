package errors

import (
	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidValueError(value interface{}, includeStack bool) RichError {
	msg := "an invalid value was found."
	err := NewRichError(codes.ErrCodeInvalidValue, msg, includeStack).AddMetaData("value", value)
	return err
}
