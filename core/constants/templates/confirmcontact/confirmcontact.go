package confirmcontact

import "github.com/calvine/goauth/core/constants/templates"

type ConfirmContactTextTemplateInput struct {
	ServiceName string
	ConfirmLink string
}

const (
	ConfirmContactTextEmail         templates.TemplateName = "confirmcontact-text-email"
	ConfirmContactTextEmailTemplate string                 = `Thank you for registering this email with {{ .ServiceName }}. Please go to the link below to confirm your contact:
	
	{{ .ConfirmLink }}`

	ConfirmContactHtmlEmail templates.TemplateName = "confirmcontact-html-email"
)
