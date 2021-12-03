package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodePasswordResetContactNotPrimary contact for password reset is not a primary contact
const ErrCodePasswordResetContactNotPrimary = "PasswordResetContactNotPrimary"

// NewPasswordResetContactNotPrimaryError creates a new specific error
func NewPasswordResetContactNotPrimaryError(contactId string, principal string, principalType string, includeStack bool) errors.RichError {
	msg := "contact for password reset is not a primary contact"
	err := errors.NewRichError(ErrCodePasswordResetContactNotPrimary, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType).WithTags([]string{"login", "security"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsPasswordResetContactNotPrimaryError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodePasswordResetContactNotPrimary
}
