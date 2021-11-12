package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeLoginPrimaryContactNotConfirmed primary contact for login is not confirmed
const ErrCodeLoginPrimaryContactNotConfirmed = "LoginPrimaryContactNotConfirmed"

// NewLoginPrimaryContactNotConfirmedError creates a new specific error
func NewLoginPrimaryContactNotConfirmedError(contactId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "primary contact for login is not confirmed"
	err := errors.NewRichError(ErrCodeLoginPrimaryContactNotConfirmed, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsLoginPrimaryContactNotConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeLoginPrimaryContactNotConfirmed
}
