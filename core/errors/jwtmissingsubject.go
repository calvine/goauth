package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTMissingSubject jwt is missing subject
const ErrCodeJWTMissingSubject = "JWTMissingSubject"

// NewJWTMissingSubjectError creates a new specific error
func NewJWTMissingSubjectError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt is missing subject"
	err := errors.NewRichError(ErrCodeJWTMissingSubject, msg).WithMetaData(fields).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTMissingSubjectError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTMissingSubject
}
