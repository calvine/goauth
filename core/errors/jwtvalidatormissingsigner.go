package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorMissingSigner jwt validator is missing a signer for signature validation
const ErrCodeJWTValidatorMissingSigner = "JWTValidatorMissingSigner"

// NewJWTValidatorMissingSignerError creates a new specific error
func NewJWTValidatorMissingSignerError(includeStack bool) errors.RichError {
	msg := "jwt validator is missing a signer for signature validation"
	err := errors.NewRichError(ErrCodeJWTValidatorMissingSigner, msg).WithTags([]string{"jwt_signer"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorMissingSignerError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorMissingSigner
}
