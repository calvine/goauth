package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"

	"time"
)

// ErrCodeJWTExipred jwt is expired
const ErrCodeJWTExipred = "JWTExipred"

// NewJWTExipredError creates a new specific error
func NewJWTExipredError(exp time.Time, fields map[string]interface{}, includeStack bool) errors.RichError {
	msg := "jwt is expired"
	err := errors.NewRichError(ErrCodeJWTExipred, msg).WithMetaData(fields).AddMetaData("exp", exp).WithTags([]string{"security", "jwt"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsJWTExipredError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeJWTExipred
}
