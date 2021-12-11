package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTNotBeforeMissing jwt not before is missing
const ErrCodeJWTNotBeforeMissing = "JWTNotBeforeMissing"

// NewJWTNotBeforeMissingError creates a new specific error
func NewJWTNotBeforeMissingError(includeStack bool) errors.RichError {
	msg := "jwt not before is missing"
	err := errors.NewRichError(ErrCodeJWTNotBeforeMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTNotBeforeMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTNotBeforeMissing
}
