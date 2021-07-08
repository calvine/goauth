package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidTypeError(actual string) RichError {
	msg := fmt.Sprintf("invalid type encountered: %s", actual)
	err := NewRichError(codes.ErrCodeInvalidType, msg).AddMetaData("actual", actual)
	return err
}
