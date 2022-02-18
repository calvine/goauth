package template

type Type string

const (
	Text Type = "text"
	Html Type = "html"
)

const (
	ConfirmContactTextEmailTemplateName = "confirmcontact-text-email"
	ConfirmContactTextEmailTemplate     = `Thank you for registering this email with {{ .ServiceName }}. Please go to the link below to confirm your contact:
	
	{{ .ConfirmLink }}`

	ConfirmContactHtmlEmailTemplateName = "confirmcontact-html-email"

	PasswordResetTextEmailTemplateName = "passwordreset-text-email"
	PasswordResetHtmlEmailTemplateName = "passwordreset-html-email"

	MagicLoginTextEmailTemplateName = "magiclogin-text-email"
	MagicLoginHtmlEmailTemplateName = "magiclogin-html-email"
)
