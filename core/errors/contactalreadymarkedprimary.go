package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeContactAlreadyMarkedPrimary contact is already marked as primary
const ErrCodeContactAlreadyMarkedPrimary = "ContactAlreadyMarkedPrimary"

// NewContactAlreadyMarkedPrimaryError creates a new specific error
func NewContactAlreadyMarkedPrimaryError(principal string, principalType string, includeStack bool) errors.RichError {
	msg := "contact is already marked as primary"
	err := errors.NewRichError(ErrCodeContactAlreadyMarkedPrimary, msg).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsContactAlreadyMarkedPrimaryError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeContactAlreadyMarkedPrimary
}
