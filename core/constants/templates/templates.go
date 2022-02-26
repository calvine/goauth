package templates

type TemplateName string
type Type string

const (
	Text Type = "text"
	Html Type = "html"
)

// TODO: write better email tempaltes...
const (
	PasswordResetTextEmail TemplateName = "passwordreset-text-email"
	PasswordResetHtmlEmail TemplateName = "passwordreset-html-email"

	MagicLoginTextEmail TemplateName = "magiclogin-text-email"
	MagicLoginHtmlEmail TemplateName = "magiclogin-html-email"
)
