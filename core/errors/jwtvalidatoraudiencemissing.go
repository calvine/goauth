package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorAudienceMissing jwt validator required audience but no audience was provided
const ErrCodeJWTValidatorAudienceMissing = "JWTValidatorAudienceMissing"

// NewJWTValidatorAudienceMissingError creates a new specific error
func NewJWTValidatorAudienceMissingError(includeStack bool) errors.RichError {
	msg := "jwt validator required audience but no audience was provided"
	err := errors.NewRichError(ErrCodeJWTValidatorAudienceMissing, msg).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorAudienceMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorAudienceMissing
}
