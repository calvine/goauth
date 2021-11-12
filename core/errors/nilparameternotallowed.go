package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNilParameterNotAllowed a nil value was encountered for a parameter, but not allowed
const ErrCodeNilParameterNotAllowed = "NilParameterNotAllowed"

// NewNilParameterNotAllowedError creates a new specific error
func NewNilParameterNotAllowedError(functionName string, parameterName string, includeStack bool) errors.RichError {
	msg := "a nil value was encountered for a parameter, but not allowed"
	err := errors.NewRichError(ErrCodeNilParameterNotAllowed, msg).AddMetaData("functionName", functionName).AddMetaData("parameterName", parameterName)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNilParameterNotAllowedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNilParameterNotAllowed
}
