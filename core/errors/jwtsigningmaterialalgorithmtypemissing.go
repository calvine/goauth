package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSigningMaterialAlgorithmTypeMissing algorithm type for jwt signing material is missing
const ErrCodeJWTSigningMaterialAlgorithmTypeMissing = "JWTSigningMaterialAlgorithmTypeMissing"

// NewJWTSigningMaterialAlgorithmTypeMissingError creates a new specific error
func NewJWTSigningMaterialAlgorithmTypeMissingError(includeStack bool) errors.RichError {
	msg := "algorithm type for jwt signing material is missing"
	err := errors.NewRichError(ErrCodeJWTSigningMaterialAlgorithmTypeMissing, msg)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSigningMaterialAlgorithmTypeMissingError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSigningMaterialAlgorithmTypeMissing
}
