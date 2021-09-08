package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeMalfomedToken token is not valid
const ErrCodeMalfomedToken = "MalfomedToken"

// NewMalfomedTokenError creates a new specific error
func NewMalfomedTokenError(fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "token is not valid"
	err := errors.NewRichError(ErrCodeMalfomedToken, msg).WithMetaData(fields)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsMalfomedTokenError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeMalfomedToken
}
