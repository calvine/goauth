package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSignatureInvalid jwt signature is invalid
const ErrCodeJWTSignatureInvalid = "JWTSignatureInvalid"

// NewJWTSignatureInvalidError creates a new specific error
func NewJWTSignatureInvalidError(JWT string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt signature is invalid"
	err := errors.NewRichError(ErrCodeJWTSignatureInvalid, msg).WithMetaData(fields).AddMetaData("JWT", JWT).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSignatureInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSignatureInvalid
}
