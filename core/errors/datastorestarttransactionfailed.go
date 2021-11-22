package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeDatastoreStartTransactionFailed data store failed to start transaction
const ErrCodeDatastoreStartTransactionFailed = "DatastoreStartTransactionFailed"

// NewDatastoreStartTransactionFailedError creates a new specific error
func NewDatastoreStartTransactionFailedError(cause error, includeStack bool) errors.RichError {
	msg := "data store failed to start transaction"
	err := errors.NewRichError(ErrCodeDatastoreStartTransactionFailed, msg).AddError(cause).WithTags([]string{"datastore"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsDatastoreStartTransactionFailedError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeDatastoreStartTransactionFailed
}
