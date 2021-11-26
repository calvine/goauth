package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeInvalidSMTPEmailOptions smtp options provided are invalid
const ErrCodeInvalidSMTPEmailOptions = "InvalidSMTPEmailOptions"

// NewInvalidSMTPEmailOptionsError creates a new specific error
func NewInvalidSMTPEmailOptionsError(reason string, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "smtp options provided are invalid"
	err := errors.NewRichError(ErrCodeInvalidSMTPEmailOptions, msg).WithMetaData(fields).AddMetaData("reason", reason).WithTags([]string{"email", "smtp"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsInvalidSMTPEmailOptionsError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeInvalidSMTPEmailOptions
}
