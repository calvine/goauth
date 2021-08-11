package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewLoginContactNotPrimaryError creates a new specific error
func NewLoginContactNotPrimaryError(contactId string, principal string, principalType string, includeStack bool) RichError {
	msg := "contact user for login is not a primary contact"
	err := NewRichError(codes.ErrCodeLoginContactNotPrimary, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
