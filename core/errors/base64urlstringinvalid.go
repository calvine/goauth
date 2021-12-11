package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeBase64URLStringInvalid base 64 url string is invalid
const ErrCodeBase64URLStringInvalid = "Base64URLStringInvalid"

// NewBase64URLStringInvalidError creates a new specific error
func NewBase64URLStringInvalidError(base64String string, includeStack bool) errors.RichError {
	msg := "base 64 url string is invalid"
	err := errors.NewRichError(ErrCodeBase64URLStringInvalid, msg).AddMetaData("base64String", base64String)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsBase64URLStringInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeBase64URLStringInvalid
}
