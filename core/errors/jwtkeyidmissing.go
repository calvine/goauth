package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTKeyIDMissing jwt key id is missing
const ErrCodeJWTKeyIDMissing = "JWTKeyIDMissing"

// NewJWTKeyIDMissingError creates a new specific error
func NewJWTKeyIDMissingError(includeStack bool) errors.RichError {
	msg := "jwt key id is missing"
	err := errors.NewRichError(ErrCodeJWTKeyIDMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTKeyIDMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTKeyIDMissing
}
