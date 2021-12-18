package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTStandardClaimsInvalid jwt standard claims are invalid
const ErrCodeJWTStandardClaimsInvalid = "JWTStandardClaimsInvalid"

// NewJWTStandardClaimsInvalidError creates a new specific error
func NewJWTStandardClaimsInvalidError(jwt string, includeStack bool) errors.RichError {
	msg := "jwt standard claims are invalid"
	err := errors.NewRichError(ErrCodeJWTStandardClaimsInvalid, msg).AddMetaData("jwt", jwt).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTStandardClaimsInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTStandardClaimsInvalid
}
