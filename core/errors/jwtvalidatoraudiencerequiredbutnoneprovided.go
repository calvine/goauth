package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorAudienceRequiredButNoneProvided jwt validator requires specific audiences but none were provided
const ErrCodeJWTValidatorAudienceRequiredButNoneProvided = "JWTValidatorAudienceRequiredButNoneProvided"

// NewJWTValidatorAudienceRequiredButNoneProvidedError creates a new specific error
func NewJWTValidatorAudienceRequiredButNoneProvidedError(includeStack bool) errors.RichError {
	msg := "jwt validator requires specific audiences but none were provided"
	err := errors.NewRichError(ErrCodeJWTValidatorAudienceRequiredButNoneProvided, msg).WithTags([]string{"config", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorAudienceRequiredButNoneProvidedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorAudienceRequiredButNoneProvided
}
