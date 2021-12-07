package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTMalformed encoded jwt is malformed
const ErrCodeJWTMalformed = "JWTMalformed"

// NewJWTMalformedError creates a new specific error
func NewJWTMalformedError(message string, JWT string, includeStack bool) errors.RichError {
	msg := "encoded jwt is malformed"
	err := errors.NewRichError(ErrCodeJWTMalformed, msg).AddMetaData("message", message).AddMetaData("JWT", JWT).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTMalformedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTMalformed
}
