package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewFailedToParseObjectIDError creates a new specific error
func NewFailedToParseObjectIDError(oid string, error error, includeStack bool) RichError {
	msg := "failed to parse an object id string"
	err := NewRichError(codes.ErrCodeFailedToParseObjectID, msg).AddMetaData("oid", oid).AddMetaData("error", error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
