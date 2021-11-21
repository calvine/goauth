package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeDatastoreTransactionAbortFailed data store transaction was aborted
const ErrCodeDatastoreTransactionAbortFailed = "DatastoreTransactionAbortFailed"

// NewDatastoreTransactionAbortFailedError creates a new specific error
func NewDatastoreTransactionAbortFailedError(cause error, transactionAbortError error, includeStack bool) errors.RichError {
	msg := "data store transaction was aborted"
	err := errors.NewRichError(ErrCodeDatastoreTransactionAbortFailed, msg).AddError(cause).AddError(transactionAbortError)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsDatastoreTransactionAbortFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeDatastoreTransactionAbortFailed
}
