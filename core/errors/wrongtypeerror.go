package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewWrongTypeError(actual, expected string) RichError {
	msg := fmt.Sprintf("wrong type found: actual: %s - expected: %s", expected, actual)
	err := NewRichError(codes.ErrCodeWrongType, msg).AddMetaData("actual", actual).AddMetaData("expected", expected)
	return err
}
