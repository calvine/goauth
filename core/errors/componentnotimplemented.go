package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewComponentNotImplementedError creates a new specific error
func NewComponentNotImplementedError(compoenentType string, missingType string, includeStack bool) RichError {
	msg := "component not implemented"
	err := NewRichError(codes.ErrCodeComponentNotImplemented, msg).AddMetaData("compoenentType", compoenentType).AddMetaData("missingType", missingType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
