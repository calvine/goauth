package templates

type TemplateName string
type Type string

const (
	Text Type = "text"
	Html Type = "html"
)

// TODO: write better email tempaltes...
const (
	ConfirmContactTextEmail         TemplateName = "confirmcontact-text-email"
	ConfirmContactTextEmailTemplate string       = `Thank you for registering this email with {{ .ServiceName }}. Please go to the link below to confirm your contact:
	
	{{ .ConfirmLink }}`

	ConfirmContactHtmlEmail TemplateName = "confirmcontact-html-email"

	PasswordResetTextEmail TemplateName = "passwordreset-text-email"
	PasswordResetHtmlEmail TemplateName = "passwordreset-html-email"

	MagicLoginTextEmail TemplateName = "magiclogin-text-email"
	MagicLoginHtmlEmail TemplateName = "magiclogin-html-email"
)

type ConfirmContactTextTemplateInput struct {
	ServiceName string
	ConfirmLink string
}
