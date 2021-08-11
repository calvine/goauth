package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewLoginFailedWrongPasswordError creates a new specific error
func NewLoginFailedWrongPasswordError(userID string, includeStack bool) RichError {
	msg := "login for user failed due to incorrect password"
	err := NewRichError(codes.ErrCodeLoginFailedWrongPassword, msg).AddMetaData("userID", userID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
