package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewRepoQueryFailedError creates a new specific error
func NewRepoQueryFailedError(queryError error, includeStack bool) RichError {
	msg := "repo query failed with error"
	err := NewRichError(codes.ErrCodeRepoQueryFailed, msg).AddMetaData("queryError", queryError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
