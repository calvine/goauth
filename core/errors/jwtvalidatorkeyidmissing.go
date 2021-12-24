package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorKeyIDMissing jwt validator required jwy key id but no jwt key id was provided
const ErrCodeJWTValidatorKeyIDMissing = "JWTValidatorKeyIDMissing"

// NewJWTValidatorKeyIDMissingError creates a new specific error
func NewJWTValidatorKeyIDMissingError(includeStack bool) errors.RichError {
	msg := "jwt validator required jwy key id but no jwt key id was provided"
	err := errors.NewRichError(ErrCodeJWTValidatorKeyIDMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorKeyIDMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorKeyIDMissing
}
