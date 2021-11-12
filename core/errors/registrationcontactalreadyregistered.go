package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeRegistrationContactAlreadyRegistered contact provided is not confirmed
const ErrCodeRegistrationContactAlreadyRegistered = "RegistrationContactAlreadyRegistered"

// NewRegistrationContactAlreadyRegisteredError creates a new specific error
func NewRegistrationContactAlreadyRegisteredError(principal, principalType, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "contact provided is not confirmed"
	err := errors.NewRichError(ErrCodeRegistrationContactAlreadyRegistered, msg).WithMetaData(fields).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsRegistrationContactAlreadyRegisteredError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeRegistrationContactAlreadyRegistered
}
