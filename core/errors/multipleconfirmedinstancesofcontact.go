package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeMultipleConfirmedInstancesOfContact multiple confirmed instance of contact found
const ErrCodeMultipleConfirmedInstancesOfContact = "MultipleConfirmedInstancesOfContact"

// NewMultipleConfirmedInstancesOfContactError creates a new specific error
func NewMultipleConfirmedInstancesOfContactError(principal string, principalType string, numOccurances int64, includeStack bool) errors.RichError {
	msg := "multiple confirmed instance of contact found"
	err := errors.NewRichError(ErrCodeMultipleConfirmedInstancesOfContact, msg).AddMetaData("principal", principal).AddMetaData("principalType", principalType).AddMetaData("numOccurances", numOccurances)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsMultipleConfirmedInstancesOfContactError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeMultipleConfirmedInstancesOfContact
}
