package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeDatastoreCommitTransactionFailed data store failed to commit transaction
const ErrCodeDatastoreCommitTransactionFailed = "DatastoreCommitTransactionFailed"

// NewDatastoreCommitTransactionFailedError creates a new specific error
func NewDatastoreCommitTransactionFailedError(cause error, includeStack bool) errors.RichError {
	msg := "data store failed to commit transaction"
	err := errors.NewRichError(ErrCodeDatastoreCommitTransactionFailed, msg).AddError(cause)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsDatastoreCommitTransactionFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeDatastoreCommitTransactionFailed
}
