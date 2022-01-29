package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeJWTSigningMaterialAlgorithmTypeNotSupported algorithm type for jwt signing material is not supported
const ErrCodeJWTSigningMaterialAlgorithmTypeNotSupported = "JWTSigningMaterialAlgorithmTypeNotSupported"

// NewJWTSigningMaterialAlgorithmTypeNotSupportedError creates a new specific error
func NewJWTSigningMaterialAlgorithmTypeNotSupportedError(algorithmType string, includeStack bool) errors.RichError {
	msg := "algorithm type for jwt signing material is not supported"
	err := errors.NewRichError(ErrCodeJWTSigningMaterialAlgorithmTypeNotSupported, msg).AddMetaData("algorithmType", algorithmType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTSigningMaterialAlgorithmTypeNotSupportedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTSigningMaterialAlgorithmTypeNotSupported
}
