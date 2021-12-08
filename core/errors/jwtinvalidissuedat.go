package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeJWTInvalidIssuedAt jwt issued at is invalid
const ErrCodeJWTInvalidIssuedAt = "JWTInvalidIssuedAt"

// NewJWTInvalidIssuedAtError creates a new specific error
func NewJWTInvalidIssuedAtError(iat time.Time, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt issued at is invalid"
	err := errors.NewRichError(ErrCodeJWTInvalidIssuedAt, msg).WithMetaData(fields).AddMetaData("iat", iat).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTInvalidIssuedAtError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTInvalidIssuedAt
}
