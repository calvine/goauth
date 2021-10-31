package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidContactType the contact type provided is not valid
const ErrCodeInvalidContactType = "InvalidContactType"

// NewInvalidContactTypeError creates a new specific error
func NewInvalidContactTypeError(principal string, principalType string, includeStack bool) errors.RichError {
	msg := "the contact type provided is not valid"
	err := errors.NewRichError(ErrCodeInvalidContactType, msg).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidContactTypeError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidContactType
}
