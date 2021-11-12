package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeExistingContactNotConfirmed existing contact is not confirmed
const ErrCodeExistingContactNotConfirmed = "ExistingContactNotConfirmed"

// NewExistingContactNotConfirmedError creates a new specific error
func NewExistingContactNotConfirmedError(contactId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "existing contact is not confirmed"
	err := errors.NewRichError(ErrCodeExistingContactNotConfirmed, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsExistingContactNotConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeExistingContactNotConfirmed
}
