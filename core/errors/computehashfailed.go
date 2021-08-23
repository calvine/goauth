package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeComputeHashFailed a failure occurred while attempting to calculate a hash
const ErrCodeComputeHashFailed = "ComputeHashFailed"

// NewComputeHashFailedError creates a new specific error
func NewComputeHashFailedError(algorithm string, error error, includeStack bool) errors.RichError {
	msg := "a failure occurred while attempting to calculate a hash"
	err := errors.NewRichError(ErrCodeComputeHashFailed, msg).AddMetaData("algorithm", algorithm).AddError(error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsComputeHashFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeComputeHashFailed
}
