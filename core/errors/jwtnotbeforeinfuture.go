package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeJWTNotBeforeInFuture jwt not before has not yet happened yet
const ErrCodeJWTNotBeforeInFuture = "JWTNotBeforeInFuture"

// NewJWTNotBeforeInFutureError creates a new specific error
func NewJWTNotBeforeInFutureError(nbf time.Time, includeStack bool) errors.RichError {
	msg := "jwt not before has not yet happened yet"
	err := errors.NewRichError(ErrCodeJWTNotBeforeInFuture, msg).AddMetaData("nbf", nbf).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTNotBeforeInFutureError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTNotBeforeInFuture
}
