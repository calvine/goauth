package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoMagicLoginTokenFound the magic login token was not found in the request
const ErrCodeNoMagicLoginTokenFound = "NoMagicLoginTokenFound"

// NewNoMagicLoginTokenFoundError creates a new specific error
func NewNoMagicLoginTokenFoundError(includeStack bool) errors.RichError {
	msg := "the magic login token was not found in the request"
	err := errors.NewRichError(ErrCodeNoMagicLoginTokenFound, msg).WithTags([]string{"login", "security"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoMagicLoginTokenFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoMagicLoginTokenFound
}
