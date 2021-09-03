package email

type EmailRecipientInfo struct {
	Bcc []string `json:"bcc"`
	Cc  []string `json:"cc"`
	To  []string `json:"to"`
}

type EmailContentInfo struct {
	Body       string             `json:"body"`
	IsHTMLBody bool               `json:"isHtmlBody"`
	Subject    string             `json:"subject"`
	Recipients EmailRecipientInfo `json:",inline"`
}
