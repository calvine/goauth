package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSigningMaterialExpired jwt signing material has expired
const ErrCodeJWTSigningMaterialExpired = "JWTSigningMaterialExpired"

// NewJWTSigningMaterialExpiredError creates a new specific error
func NewJWTSigningMaterialExpiredError(keyID string, includeStack bool) errors.RichError {
	msg := "jwt signing material has expired"
	err := errors.NewRichError(ErrCodeJWTSigningMaterialExpired, msg).AddMetaData("keyID", keyID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSigningMaterialExpiredError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSigningMaterialExpired
}
