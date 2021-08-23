package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeBcryptPasswordHashError an error occurred while processing bcrypt hash for password
const ErrCodeBcryptPasswordHashError = "BcryptPasswordHashError"

// NewBcryptPasswordHashErrorError creates a new specific error
func NewBcryptPasswordHashErrorError(assetId string, error error, includeStack bool) errors.RichError {
	msg := "an error occurred while processing bcrypt hash for password"
	err := errors.NewRichError(ErrCodeBcryptPasswordHashError, msg).AddMetaData("assetId", assetId).AddError(error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsBcryptPasswordHashErrorError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeBcryptPasswordHashError
}
