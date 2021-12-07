package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTAlgorithmInvalid jwt algorithm is invalid
const ErrCodeJWTAlgorithmInvalid = "JWTAlgorithmInvalid"

// NewJWTAlgorithmInvalidError creates a new specific error
func NewJWTAlgorithmInvalidError(providedAlgorithm string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt algorithm is invalid"
	err := errors.NewRichError(ErrCodeJWTAlgorithmInvalid, msg).WithMetaData(fields).AddMetaData("providedAlgorithm", providedAlgorithm).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTAlgorithmInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTAlgorithmInvalid
}
