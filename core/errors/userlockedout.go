package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeUserLockedOut attempted login by locked out user
const ErrCodeUserLockedOut = "UserLockedOut"

// NewUserLockedOutError creates a new specific error
func NewUserLockedOutError(userID string, includeStack bool) errors.RichError {
	msg := "attempted login by locked out user"
	err := errors.NewRichError(ErrCodeUserLockedOut, msg).AddMetaData("userID", userID).WithTags([]string{"login", "security"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsUserLockedOutError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeUserLockedOut
}
