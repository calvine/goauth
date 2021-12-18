package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorAllowAnyIssuerAndExpectedIssuerProvided jwt validator cannot have allow any issuer set to true and expected issuer have a value
const ErrCodeJWTValidatorAllowAnyIssuerAndExpectedIssuerProvided = "JWTValidatorAllowAnyIssuerAndExpectedIssuerProvided"

// NewJWTValidatorAllowAnyIssuerAndExpectedIssuerProvidedError creates a new specific error
func NewJWTValidatorAllowAnyIssuerAndExpectedIssuerProvidedError(includeStack bool) errors.RichError {
	msg := "jwt validator cannot have allow any issuer set to true and expected issuer have a value"
	err := errors.NewRichError(ErrCodeJWTValidatorAllowAnyIssuerAndExpectedIssuerProvided, msg).WithTags([]string{"config", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorAllowAnyIssuerAndExpectedIssuerProvidedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorAllowAnyIssuerAndExpectedIssuerProvided
}
