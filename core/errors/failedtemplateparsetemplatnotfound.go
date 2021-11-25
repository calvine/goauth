package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeFailedTemplateParseTemplatNotFound template not found in fs with path provided
const ErrCodeFailedTemplateParseTemplatNotFound = "FailedTemplateParseTemplatNotFound"

// NewFailedTemplateParseTemplatNotFoundError creates a new specific error
func NewFailedTemplateParseTemplatNotFoundError(templateName string, templatePath string, parseError error, includeStack bool) errors.RichError {
	msg := "template not found in fs with path provided"
	err := errors.NewRichError(ErrCodeFailedTemplateParseTemplatNotFound, msg).AddMetaData("templateName", templateName).AddMetaData("templatePath", templatePath).AddError(parseError).WithTags([]string{"template"})
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsFailedTemplateParseTemplatNotFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeFailedTemplateParseTemplatNotFound
}
