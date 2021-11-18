package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeDatastoreTransactionAborted data store transaction was aborted
const ErrCodeDatastoreTransactionAborted = "DatastoreTransactionAborted"

// NewDatastoreTransactionAbortedError creates a new specific error
func NewDatastoreTransactionAbortedError(cause error, includeStack bool) errors.RichError {
	msg := "data store transaction was aborted"
	err := errors.NewRichError(ErrCodeDatastoreTransactionAborted, msg).AddError(cause)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsDatastoreTransactionAbortedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeDatastoreTransactionAborted
}
