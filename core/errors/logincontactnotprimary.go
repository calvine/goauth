package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeLoginContactNotPrimary contact for login is not a primary contact
const ErrCodeLoginContactNotPrimary = "LoginContactNotPrimary"

// NewLoginContactNotPrimaryError creates a new specific error
func NewLoginContactNotPrimaryError(contactId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "contact for login is not a primary contact"
	err := errors.NewRichError(ErrCodeLoginContactNotPrimary, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsLoginContactNotPrimaryError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeLoginContactNotPrimary
}
