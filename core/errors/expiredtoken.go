package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeExpiredToken token has expired
const ErrCodeExpiredToken = "ExpiredToken"

// NewExpiredTokenError creates a new specific error
func NewExpiredTokenError(id string, tokenType string, expiredOn time.Time, includeStack bool) errors.RichError {
	msg := "token has expired"
	err := errors.NewRichError(ErrCodeExpiredToken, msg).AddMetaData("id", id).AddMetaData("tokenType", tokenType).AddMetaData("expiredOn", expiredOn)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsExpiredTokenError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeExpiredToken
}
