package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTExpireMissing jwt expire is missing
const ErrCodeJWTExpireMissing = "JWTExpireMissing"

// NewJWTExpireMissingError creates a new specific error
func NewJWTExpireMissingError(includeStack bool) errors.RichError {
	msg := "jwt expire is missing"
	err := errors.NewRichError(ErrCodeJWTExpireMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTExpireMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTExpireMissing
}
