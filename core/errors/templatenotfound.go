package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeTemplateNotFound template not found in data store
const ErrCodeTemplateNotFound = "TemplateNotFound"

// NewTemplateNotFoundError creates a new specific error
func NewTemplateNotFoundError(templateName string, includeStack bool) errors.RichError {
	msg := "template not found in data store"
	err := errors.NewRichError(ErrCodeTemplateNotFound, msg).AddMetaData("templateName", templateName).WithTags([]string{"template"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsTemplateNotFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeTemplateNotFound
}
