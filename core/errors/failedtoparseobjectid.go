package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedToParseObjectID failed to parse an object id string
const ErrCodeFailedToParseObjectID = "FailedToParseObjectID"

// NewFailedToParseObjectIDError creates a new specific error
func NewFailedToParseObjectIDError(oid string, error error, includeStack bool) errors.RichError {
	msg := "failed to parse an object id string"
	err := errors.NewRichError(ErrCodeFailedToParseObjectID, msg).AddMetaData("oid", oid).AddError(error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedToParseObjectIDError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedToParseObjectID
}
