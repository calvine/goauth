package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeRegisteredContactNotConfirmed existing contact is not confirmed
const ErrCodeRegisteredContactNotConfirmed = "RegisteredContactNotConfirmed"

// NewRegisteredContactNotConfirmedError creates a new specific error
func NewRegisteredContactNotConfirmedError(contactId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "existing contact is not confirmed"
	err := errors.NewRichError(ErrCodeRegisteredContactNotConfirmed, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsRegisteredContactNotConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeRegisteredContactNotConfirmed
}
