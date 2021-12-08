package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorNoAlgorithmSpecified jwt validator has no algorithms specified
const ErrCodeJWTValidatorNoAlgorithmSpecified = "JWTValidatorNoAlgorithmSpecified"

// NewJWTValidatorNoAlgorithmSpecifiedError creates a new specific error
func NewJWTValidatorNoAlgorithmSpecifiedError(includeStack bool) errors.RichError {
	msg := "jwt validator has no algorithms specified"
	err := errors.NewRichError(ErrCodeJWTValidatorNoAlgorithmSpecified, msg).WithTags([]string{"config", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorNoAlgorithmSpecifiedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorNoAlgorithmSpecified
}
