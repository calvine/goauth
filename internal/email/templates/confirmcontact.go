package templates

type ConfirmContactEmailTemplateParams struct {
	ServiceName string `json:"serviceName"`
	ConfirmLink string `json:"confirmLink"`
}

// ConfirmContactEmailTemplate_PlainText is a plain text, so it may not allow for clicking the link
// FIXME: make an html template...
const ConfirmContactEmailTemplate_PlainText = `
Thank you for registering an account with {{ .ServiceName }}! 
Please use the link below to confirm your account.

{{ .ConfirmLink }}
`
