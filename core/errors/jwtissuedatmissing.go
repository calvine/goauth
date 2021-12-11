package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTIssuedAtMissing jwt issued at is missing
const ErrCodeJWTIssuedAtMissing = "JWTIssuedAtMissing"

// NewJWTIssuedAtMissingError creates a new specific error
func NewJWTIssuedAtMissingError(includeStack bool) errors.RichError {
	msg := "jwt issued at is missing"
	err := errors.NewRichError(ErrCodeJWTIssuedAtMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTIssuedAtMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTIssuedAtMissing
}
