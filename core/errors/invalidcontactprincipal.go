package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidContactPrincipal the contact principal provided is not valid
const ErrCodeInvalidContactPrincipal = "InvalidContactPrincipal"

// NewInvalidContactPrincipalError creates a new specific error
func NewInvalidContactPrincipalError(principal string, principalType string, includeStack bool) errors.RichError {
	msg := "the contact principal provided is not valid"
	err := errors.NewRichError(ErrCodeInvalidContactPrincipal, msg).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidContactPrincipalError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidContactPrincipal
}
