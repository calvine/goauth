package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTIDMissing jwt id (jti) is missing
const ErrCodeJWTIDMissing = "JWTIDMissing"

// NewJWTIDMissingError creates a new specific error
func NewJWTIDMissingError(includeStack bool) errors.RichError {
	msg := "jwt id (jti) is missing"
	err := errors.NewRichError(ErrCodeJWTIDMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTIDMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTIDMissing
}
