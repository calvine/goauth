package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeRegistrationContactAlreadyConfirmed contact for registration is already confirmed
const ErrCodeRegistrationContactAlreadyConfirmed = "RegistrationContactAlreadyConfirmed"

// NewRegistrationContactAlreadyConfirmedError creates a new specific error
func NewRegistrationContactAlreadyConfirmedError(principal string, principalType string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "contact for registration is already confirmed"
	err := errors.NewRichError(ErrCodeRegistrationContactAlreadyConfirmed, msg).WithMetaData(fields).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsRegistrationContactAlreadyConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeRegistrationContactAlreadyConfirmed
}
