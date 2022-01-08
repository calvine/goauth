package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedToDecodeRepoData decoding data from repo call failed
const ErrCodeFailedToDecodeRepoData = "FailedToDecodeRepoData"

// NewFailedToDecodeRepoDataError creates a new specific error
func NewFailedToDecodeRepoDataError(decodeTargetType string, cause error, includeStack bool) errors.RichError {
	msg := "decoding data from repo call failed"
	err := errors.NewRichError(ErrCodeFailedToDecodeRepoData, msg).AddMetaData("decodeTargetType", decodeTargetType).AddError(cause).WithTags([]string{"repo"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedToDecodeRepoDataError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedToDecodeRepoData
}
