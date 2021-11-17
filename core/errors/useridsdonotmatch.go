package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeUserIDsDoNotMatch user id provided do not match
const ErrCodeUserIDsDoNotMatch = "UserIDsDoNotMatch"

// NewUserIDsDoNotMatchError creates a new specific error
func NewUserIDsDoNotMatchError(userID1 string, userID2 string, includeStack bool) errors.RichError {
	msg := "user id provided do not match"
	err := errors.NewRichError(ErrCodeUserIDsDoNotMatch, msg).AddMetaData("userID1", userID1).AddMetaData("userID2", userID2)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsUserIDsDoNotMatchError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeUserIDsDoNotMatch
}
