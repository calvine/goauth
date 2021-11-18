package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeLoginFailedWrongPassword login for user failed due to incorrect password
const ErrCodeLoginFailedWrongPassword = "LoginFailedWrongPassword"

// NewLoginFailedWrongPasswordError creates a new specific error
func NewLoginFailedWrongPasswordError(userID string, includeStack bool) errors.RichError {
	msg := "login for user failed due to incorrect password"
	err := errors.NewRichError(ErrCodeLoginFailedWrongPassword, msg).AddMetaData("userID", userID).WithTags([]string{"security"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsLoginFailedWrongPasswordError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeLoginFailedWrongPassword
}
