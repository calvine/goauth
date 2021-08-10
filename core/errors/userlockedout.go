package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewUserLockedOutError creates a new specific error
func NewUserLockedOutError(userID string, includeStack bool) RichError {
	msg := "attempted login by locked out user"
	err := NewRichError(codes.ErrCodeUserLockedOut, msg).AddMetaData("userID", userID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
