package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"time"

	"github.com/calvine/richerror/errors"
)

// ErrCodeExpiredPasswordExpirationToken password expiration token has expired
const ErrCodeExpiredPasswordExpirationToken = "ExpiredPasswordExpirationToken"

// NewExpiredPasswordExpirationTokenError creates a new specific error
func NewExpiredPasswordExpirationTokenError(passwordResetToken string, passwordResetTokenExpiration time.Time, includeStack bool) errors.RichError {
	msg := "password expiration token has expired"
	err := errors.NewRichError(ErrCodeExpiredPasswordExpirationToken, msg).AddMetaData("passwordResetToken", passwordResetToken).AddMetaData("passwordResetTokenExpiration", passwordResetTokenExpiration)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsExpiredPasswordExpirationTokenError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeExpiredPasswordExpirationToken
}
