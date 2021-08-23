package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeWrongType unexpected type encountered
const ErrCodeWrongType = "WrongType"

// NewWrongTypeError creates a new specific error
func NewWrongTypeError(actual string, expected string, includeStack bool) errors.RichError {
	msg := "unexpected type encountered"
	err := errors.NewRichError(ErrCodeWrongType, msg).AddMetaData("actual", actual).AddMetaData("expected", expected)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsWrongTypeError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeWrongType
}
