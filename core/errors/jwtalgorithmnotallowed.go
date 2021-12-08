package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTAlgorithmNotAllowed jwt algorithm is not allowed
const ErrCodeJWTAlgorithmNotAllowed = "JWTAlgorithmNotAllowed"

// NewJWTAlgorithmNotAllowedError creates a new specific error
func NewJWTAlgorithmNotAllowedError(algorithm string, includeStack bool) errors.RichError {
	msg := "jwt algorithm is not allowed"
	err := errors.NewRichError(ErrCodeJWTAlgorithmNotAllowed, msg).AddMetaData("algorithm", algorithm).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTAlgorithmNotAllowedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTAlgorithmNotAllowed
}
