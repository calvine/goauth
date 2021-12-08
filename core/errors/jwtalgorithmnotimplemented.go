package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTAlgorithmNotImplemented jwt algorithm is not implemented
const ErrCodeJWTAlgorithmNotImplemented = "JWTAlgorithmNotImplemented"

// NewJWTAlgorithmNotImplementedError creates a new specific error
func NewJWTAlgorithmNotImplementedError(algorithm string, includeStack bool) errors.RichError {
	msg := "jwt algorithm is not implemented"
	err := errors.NewRichError(ErrCodeJWTAlgorithmNotImplemented, msg).AddMetaData("algorithm", algorithm).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTAlgorithmNotImplementedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTAlgorithmNotImplemented
}
