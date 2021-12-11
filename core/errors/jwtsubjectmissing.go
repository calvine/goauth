package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSubjectMissing jwt subject is missing
const ErrCodeJWTSubjectMissing = "JWTSubjectMissing"

// NewJWTSubjectMissingError creates a new specific error
func NewJWTSubjectMissingError(includeStack bool) errors.RichError {
	msg := "jwt subject is missing"
	err := errors.NewRichError(ErrCodeJWTSubjectMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSubjectMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSubjectMissing
}
