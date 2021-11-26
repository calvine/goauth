package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedToSendEmail failed to send email message
const ErrCodeFailedToSendEmail = "FailedToSendEmail"

// NewFailedToSendEmailError creates a new specific error
func NewFailedToSendEmailError(cause error, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "failed to send email message"
	err := errors.NewRichError(ErrCodeFailedToSendEmail, msg).WithMetaData(fields).AddError(cause).WithTags([]string{"email"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedToSendEmailError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedToSendEmail
}
