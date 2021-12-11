package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTIssuerMissing jwt issuer is missing
const ErrCodeJWTIssuerMissing = "JWTIssuerMissing"

// NewJWTIssuerMissingError creates a new specific error
func NewJWTIssuerMissingError(includeStack bool) errors.RichError {
	msg := "jwt issuer is missing"
	err := errors.NewRichError(ErrCodeJWTIssuerMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTIssuerMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTIssuerMissing
}
