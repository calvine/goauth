package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeWrongTokenType token type does not match expected token type
const ErrCodeWrongTokenType = "WrongTokenType"

// NewWrongTokenTypeError creates a new specific error
func NewWrongTokenTypeError(id string, tokenType string, expectedTokenType string, includeStack bool) errors.RichError {
	msg := "token type does not match expected token type"
	err := errors.NewRichError(ErrCodeWrongTokenType, msg).AddMetaData("id", id).AddMetaData("tokenType", tokenType).AddMetaData("expectedTokenType", expectedTokenType).WithTags([]string{"token"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsWrongTokenTypeError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeWrongTokenType
}
