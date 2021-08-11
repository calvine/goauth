package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewComputeHashFailedError creates a new specific error
func NewComputeHashFailedError(algorithm string, error error, includeStack bool) RichError {
	msg := "a failure occurred while attempting to calculate a hash"
	err := NewRichError(codes.ErrCodeComputeHashFailed, msg).AddMetaData("algorithm", algorithm).AddMetaData("error", error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
