package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeMagicLoginTokenNoUserID the magic login token was not tied to a user id
const ErrCodeMagicLoginTokenNoUserID = "MagicLoginTokenNoUserID"

// NewMagicLoginTokenNoUserIDError creates a new specific error
func NewMagicLoginTokenNoUserIDError(tokenString string, includeStack bool) errors.RichError {
	msg := "the magic login token was not tied to a user id"
	err := errors.NewRichError(ErrCodeMagicLoginTokenNoUserID, msg).AddMetaData("tokenString", tokenString).WithTags([]string{"login", "security"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsMagicLoginTokenNoUserIDError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeMagicLoginTokenNoUserID
}
