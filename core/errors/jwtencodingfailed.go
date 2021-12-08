package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTEncodingFailed jwt encoding failed
const ErrCodeJWTEncodingFailed = "JWTEncodingFailed"

// NewJWTEncodingFailedError creates a new specific error
func NewJWTEncodingFailedError(cause error, includeStack bool) errors.RichError {
	msg := "jwt encoding failed"
	err := errors.NewRichError(ErrCodeJWTEncodingFailed, msg).AddError(cause).WithTags([]string{"jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTEncodingFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTEncodingFailed
}
