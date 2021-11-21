package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeDatastoreTransactionFailed data store transaction was aborted
const ErrCodeDatastoreTransactionFailed = "DatastoreTransactionFailed"

// NewDatastoreTransactionFailedError creates a new specific error
func NewDatastoreTransactionFailedError(cause error, includeStack bool) errors.RichError {
	msg := "data store transaction was aborted"
	err := errors.NewRichError(ErrCodeDatastoreTransactionFailed, msg).AddError(cause)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsDatastoreTransactionFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeDatastoreTransactionFailed
}
