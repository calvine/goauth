package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeBase64DecodeStringFailed base 64 string is invalid
const ErrCodeBase64DecodeStringFailed = "Base64DecodeStringFailed"

// NewBase64DecodeStringFailedError creates a new specific error
func NewBase64DecodeStringFailedError(cause error, base64String string, includeStack bool) errors.RichError {
	msg := "base 64 string is invalid"
	err := errors.NewRichError(ErrCodeBase64DecodeStringFailed, msg).AddError(cause).AddMetaData("base64String", base64String)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsBase64DecodeStringFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeBase64DecodeStringFailed
}
