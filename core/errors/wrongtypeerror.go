package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewWrongTypeError(actual, expected string, includeStack bool) RichError {
	msg := fmt.Sprintf("wrong type found: actual: %s - expected: %s", expected, actual)
	err := NewRichError(codes.ErrCodeWrongType, msg).AddMetaData("actual", actual).AddMetaData("expected", expected)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
