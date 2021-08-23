package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidType invalid type encountered
const ErrCodeInvalidType = "InvalidType"

// NewInvalidTypeError creates a new specific error
func NewInvalidTypeError(typeEncountered string, includeStack bool) errors.RichError {
	msg := "invalid type encountered"
	err := errors.NewRichError(ErrCodeInvalidType, msg).AddMetaData("typeEncountered", typeEncountered)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidTypeError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidType
}
