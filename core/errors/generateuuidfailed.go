package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeGenerateUUIDFailed failed to generate uuid
const ErrCodeGenerateUUIDFailed = "GenerateUUIDFailed"

// NewGenerateUUIDFailedError creates a new specific error
func NewGenerateUUIDFailedError(generationError error, includeStack bool) errors.RichError {
	msg := "failed to generate uuid"
	err := errors.NewRichError(ErrCodeGenerateUUIDFailed, msg).AddError(generationError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsGenerateUUIDFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeGenerateUUIDFailed
}
