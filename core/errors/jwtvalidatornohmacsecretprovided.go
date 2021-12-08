package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorNoHMACSecretProvided jwt validator allows HMAC algorithms but was not provided an HMAC secret
const ErrCodeJWTValidatorNoHMACSecretProvided = "JWTValidatorNoHMACSecretProvided"

// NewJWTValidatorNoHMACSecretProvidedError creates a new specific error
func NewJWTValidatorNoHMACSecretProvidedError(includeStack bool) errors.RichError {
	msg := "jwt validator allows HMAC algorithms but was not provided an HMAC secret"
	err := errors.NewRichError(ErrCodeJWTValidatorNoHMACSecretProvided, msg).WithTags([]string{"config", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorNoHMACSecretProvidedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorNoHMACSecretProvided
}
