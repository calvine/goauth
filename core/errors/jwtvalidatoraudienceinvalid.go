package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTValidatorAudienceInvalid jwt audience not allowed
const ErrCodeJWTValidatorAudienceInvalid = "JWTValidatorAudienceInvalid"

// NewJWTValidatorAudienceInvalidError creates a new specific error
func NewJWTValidatorAudienceInvalidError(audience string, includeStack bool) errors.RichError {
	msg := "jwt audience not allowed"
	err := errors.NewRichError(ErrCodeJWTValidatorAudienceInvalid, msg).AddMetaData("audience", audience).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTValidatorAudienceInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTValidatorAudienceInvalid
}
