package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorAudienceRequiredButNoneProvided jwt validator required audience but no allowed audience
const ErrCodeJWTValidatorAudienceRequiredButNoneProvided = "JWTValidatorAudienceRequiredButNoneProvided"

// NewJWTValidatorAudienceRequiredButNoneProvidedError creates a new specific error
func NewJWTValidatorAudienceRequiredButNoneProvidedError(includeStack bool) errors.RichError {
	msg := "jwt validator required audience but no allowed audience"
	err := errors.NewRichError(ErrCodeJWTValidatorAudienceRequiredButNoneProvided, msg).WithTags([]string{"config", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorAudienceRequiredButNoneProvidedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorAudienceRequiredButNoneProvided
}
