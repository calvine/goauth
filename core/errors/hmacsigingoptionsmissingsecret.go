package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeHMACSigingOptionsMissingSecret hmac signing options secret was empty
const ErrCodeHMACSigingOptionsMissingSecret = "HMACSigingOptionsMissingSecret"

// NewHMACSigingOptionsMissingSecretError creates a new specific error
func NewHMACSigingOptionsMissingSecretError(includeStack bool) errors.RichError {
	msg := "hmac signing options secret was empty"
	err := errors.NewRichError(ErrCodeHMACSigingOptionsMissingSecret, msg).WithTags([]string{"jwt_signer"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsHMACSigingOptionsMissingSecretError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeHMACSigingOptionsMissingSecret
}
