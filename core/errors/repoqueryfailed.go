package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeRepoQueryFailed repo query failed with error
const ErrCodeRepoQueryFailed = "RepoQueryFailed"

// NewRepoQueryFailedError creates a new specific error
func NewRepoQueryFailedError(queryError error, includeStack bool) errors.RichError {
	msg := "repo query failed with error"
	err := errors.NewRichError(ErrCodeRepoQueryFailed, msg).AddError(queryError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsRepoQueryFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeRepoQueryFailed
}
