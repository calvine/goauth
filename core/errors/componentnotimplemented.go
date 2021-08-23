package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeComponentNotImplemented component not implemented
const ErrCodeComponentNotImplemented = "ComponentNotImplemented"

// NewComponentNotImplementedError creates a new specific error
func NewComponentNotImplementedError(compoenentType string, missingType string, includeStack bool) errors.RichError {
	msg := "component not implemented"
	err := errors.NewRichError(ErrCodeComponentNotImplemented, msg).AddMetaData("compoenentType", compoenentType).AddMetaData("missingType", missingType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsComponentNotImplementedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeComponentNotImplemented
}
