package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeJWTIssuedAtInvalid jwt issued at is invalid
const ErrCodeJWTIssuedAtInvalid = "JWTIssuedAtInvalid"

// NewJWTIssuedAtInvalidError creates a new specific error
func NewJWTIssuedAtInvalidError(iat time.Time, includeStack bool) errors.RichError {
	msg := "jwt issued at is invalid"
	err := errors.NewRichError(ErrCodeJWTIssuedAtInvalid, msg).AddMetaData("iat", iat).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTIssuedAtInvalidError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTIssuedAtInvalid
}
