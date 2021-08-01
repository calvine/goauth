package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewUserLockedOutError(userID string, includeStack bool) RichError {
	msg := fmt.Sprintf("attempted login by locked out user with id: %s", userID)
	err := NewRichError(codes.ErrCodeUserLockedOut, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
