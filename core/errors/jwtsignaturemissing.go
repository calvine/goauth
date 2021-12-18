package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSignatureMissing jwt signature is missing
const ErrCodeJWTSignatureMissing = "JWTSignatureMissing"

// NewJWTSignatureMissingError creates a new specific error
func NewJWTSignatureMissingError(jwt string, includeStack bool) errors.RichError {
	msg := "jwt signature is missing"
	err := errors.NewRichError(ErrCodeJWTSignatureMissing, msg).AddMetaData("jwt", jwt).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSignatureMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSignatureMissing
}
