package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeContactAlreadyConfirmed contact is already confirmed
const ErrCodeContactAlreadyConfirmed = "ContactAlreadyConfirmed"

// NewContactAlreadyConfirmedError creates a new specific error
func NewContactAlreadyConfirmedError(userId string, contactId string, principal string, principalType string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "contact is already confirmed"
	err := errors.NewRichError(ErrCodeContactAlreadyConfirmed, msg).WithMetaData(fields).AddMetaData("userId", userId).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsContactAlreadyConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeContactAlreadyConfirmed
}
