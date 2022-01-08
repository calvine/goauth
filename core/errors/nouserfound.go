package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoUserFound no user found for given query
const ErrCodeNoUserFound = "NoUserFound"

// NewNoUserFoundError creates a new specific error
func NewNoUserFoundError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "no user found for given query"
	err := errors.NewRichError(ErrCodeNoUserFound, msg).WithMetaData(fields).WithTags([]string{"repo"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoUserFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoUserFound
}
