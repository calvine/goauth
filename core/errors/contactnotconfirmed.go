package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeContactNotConfirmed contact is not confirmed
const ErrCodeContactNotConfirmed = "ContactNotConfirmed"

// NewContactNotConfirmedError creates a new specific error
func NewContactNotConfirmedError(principal string, principalType string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "contact is not confirmed"
	err := errors.NewRichError(ErrCodeContactNotConfirmed, msg).WithMetaData(fields).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsContactNotConfirmedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeContactNotConfirmed
}
