package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewBcryptPasswordHashErrorError creates a new specific error
func NewBcryptPasswordHashErrorError(assetId string, error error, includeStack bool) RichError {
	msg := "an error occurred while processing bcrypt hash for password"
	err := NewRichError(codes.ErrCodeBcryptPasswordHashError, msg).AddMetaData("assetId", assetId).AddError(error)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
