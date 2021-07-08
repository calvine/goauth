package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewWrongTypeError(expected, actual string) RichError {
	msg := fmt.Sprintf("wrong type found: expected: %s - actual: %s", expected, actual)
	err := NewRichError(codes.ErrCodeWrongType, msg).AddMetaData("expected", expected).AddMetaData("actual", actual)
	return err
}
