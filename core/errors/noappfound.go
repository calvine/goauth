package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoAppFound no app found for given query
const ErrCodeNoAppFound = "NoAppFound"

// NewNoAppFoundError creates a new specific error
func NewNoAppFoundError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "no app found for given query"
	err := errors.NewRichError(ErrCodeNoAppFound, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoAppFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoAppFound
}
