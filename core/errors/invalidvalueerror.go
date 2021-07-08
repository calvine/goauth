package errors

import (
	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidValueError(value interface{}) RichError {
	err := NewRichError(codes.ErrCodeInvalidValue, "an invalid value was found.").AddMetaData("value", value)
	return err
}
