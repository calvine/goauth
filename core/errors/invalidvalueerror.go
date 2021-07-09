package errors

import (
	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidValueError(value interface{}, includeStack bool) RichError {
	msg := "an invalid value was found."
	err := NewRichError(codes.ErrCodeInvalidValue, msg).AddMetaData("value", value)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
