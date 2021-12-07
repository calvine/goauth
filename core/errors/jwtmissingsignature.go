package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTMissingSignature jwt is missing a signature
const ErrCodeJWTMissingSignature = "JWTMissingSignature"

// NewJWTMissingSignatureError creates a new specific error
func NewJWTMissingSignatureError(JWT string, includeStack bool) errors.RichError {
	msg := "jwt is missing a signature"
	err := errors.NewRichError(ErrCodeJWTMissingSignature, msg).AddMetaData("JWT", JWT).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTMissingSignatureError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTMissingSignature
}
