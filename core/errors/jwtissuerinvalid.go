package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTIssuerInvalid jwt issuer is invalid
const ErrCodeJWTIssuerInvalid = "JWTIssuerInvalid"

// NewJWTIssuerInvalidError creates a new specific error
func NewJWTIssuerInvalidError(actual string, expected string, includeStack bool) errors.RichError {
	msg := "jwt issuer is invalid"
	err := errors.NewRichError(ErrCodeJWTIssuerInvalid, msg).AddMetaData("actual", actual).AddMetaData("expected", expected).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTIssuerInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTIssuerInvalid
}
