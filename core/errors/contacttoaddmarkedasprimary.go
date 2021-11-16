package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeContactToAddMarkedAsPrimary cannot add contact that is marked as primary must go through primary contact reassignment flow
const ErrCodeContactToAddMarkedAsPrimary = "ContactToAddMarkedAsPrimary"

// NewContactToAddMarkedAsPrimaryError creates a new specific error
func NewContactToAddMarkedAsPrimaryError(userId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "cannot add contact that is marked as primary must go through primary contact reassignment flow"
	err := errors.NewRichError(ErrCodeContactToAddMarkedAsPrimary, msg).AddMetaData("userId", userId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsContactToAddMarkedAsPrimaryError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeContactToAddMarkedAsPrimary
}
