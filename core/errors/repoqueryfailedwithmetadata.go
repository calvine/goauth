package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeRepoQueryFailedWithMetaData repo query failed with error
const ErrCodeRepoQueryFailedWithMetaData = "RepoQueryFailedWithMetaData"

// NewRepoQueryFailedWithMetaDataError creates a new specific error
func NewRepoQueryFailedWithMetaDataError(queryError error, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "repo query failed with error"
	err := errors.NewRichError(ErrCodeRepoQueryFailedWithMetaData, msg).WithMetaData(fields).AddError(queryError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsRepoQueryFailedWithMetaDataError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeRepoQueryFailedWithMetaData
}
