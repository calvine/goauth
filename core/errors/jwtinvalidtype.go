package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTInvalidType jwt type is invalid
const ErrCodeJWTInvalidType = "JWTInvalidType"

// NewJWTInvalidTypeError creates a new specific error
func NewJWTInvalidTypeError(providedType string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt type is invalid"
	err := errors.NewRichError(ErrCodeJWTInvalidType, msg).WithMetaData(fields).AddMetaData("providedType", providedType).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTInvalidTypeError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTInvalidType
}
