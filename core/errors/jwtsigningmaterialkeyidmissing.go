package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSigningMaterialKeyIDMissing key id for jwt signing material is missing
const ErrCodeJWTSigningMaterialKeyIDMissing = "JWTSigningMaterialKeyIDMissing"

// NewJWTSigningMaterialKeyIDMissingError creates a new specific error
func NewJWTSigningMaterialKeyIDMissingError(includeStack bool) errors.RichError {
	msg := "key id for jwt signing material is missing"
	err := errors.NewRichError(ErrCodeJWTSigningMaterialKeyIDMissing, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSigningMaterialKeyIDMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSigningMaterialKeyIDMissing
}
