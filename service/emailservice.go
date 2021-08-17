package service

import "github.com/calvine/goauth/core/errors"

const (
	MockEmailService = "mock"
	SMTPEmailService = "smtp"
)

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

type EmailService interface {
	SendPlainTextEmail(to []string, subject, body string) errors.RichError
}

type mockEmailService struct{}

func (mse mockEmailService) SendPlainTextEmail(to []string, subject, body string) errors.RichError {
	return nil
}

func NewEmailService(serviceType string, options interface{}) (EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	// case SMTPEmailService: // TODO: implement this...
	default:
		return nil, errors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}
