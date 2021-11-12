package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeContactToAddAlreadyConfirmed contact to add to account is already confirmed
const ErrCodeContactToAddAlreadyConfirmed = "ContactToAddAlreadyConfirmed"

// NewContactToAddAlreadyConfirmedError creates a new specific error
func NewContactToAddAlreadyConfirmedError(userId string, principal string, principalType string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "contact to add to account is already confirmed"
	err := errors.NewRichError(ErrCodeContactToAddAlreadyConfirmed, msg).WithMetaData(fields).AddMetaData("userId", userId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsContactToAddAlreadyConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeContactToAddAlreadyConfirmed
}
