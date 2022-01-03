package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSigningMaterialKeyIDNotUnique key id for jwt signing material already exists
const ErrCodeJWTSigningMaterialKeyIDNotUnique = "JWTSigningMaterialKeyIDNotUnique"

// NewJWTSigningMaterialKeyIDNotUniqueError creates a new specific error
func NewJWTSigningMaterialKeyIDNotUniqueError(keyID string, includeStack bool) errors.RichError {
	msg := "key id for jwt signing material already exists"
	err := errors.NewRichError(ErrCodeJWTSigningMaterialKeyIDNotUnique, msg).AddMetaData("keyID", keyID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSigningMaterialKeyIDNotUniqueError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSigningMaterialKeyIDNotUnique
}
