package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidTypeError(actual string, includeStack bool) RichError {
	msg := fmt.Sprintf("invalid type encountered: %s", actual)
	err := NewRichError(codes.ErrCodeInvalidType, msg).AddMetaData("actual", actual)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
