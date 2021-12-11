package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeJWTExpired jwt is expired
const ErrCodeJWTExpired = "JWTExpired"

// NewJWTExpiredError creates a new specific error
func NewJWTExpiredError(exp time.Time, includeStack bool) errors.RichError {
	msg := "jwt is expired"
	err := errors.NewRichError(ErrCodeJWTExpired, msg).AddMetaData("exp", exp).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTExpiredError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTExpired
}
