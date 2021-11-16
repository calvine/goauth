package service

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	coreServices "github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

const (
	MockEmailService  = "mock"
	NoOpEmailService  = "noop"
	StackEmailService = "stack"
	SMTPEmailService  = "smtp"
)

type TestEmailMessage struct {
	To      []string
	Subject string
	Body    string
}

func NewEmailService(serviceType string, options interface{}) (coreServices.EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	case NoOpEmailService:
		return noopEmailService{}, nil
	case StackEmailService:
		return NewStackEmailService(), nil
	// case SMTPEmailService: // TODO: implement this...
	default:
		return nil, coreerrors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}

type noopEmailService struct{}

func (noopEmailService) GetName() string {
	return "noopEmailService"
}

func (ne noopEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ne.GetName(), "SendPlainTextEmail")
	defer span.End()
	return nil
}

type mockEmailService struct{}

func (mockEmailService) GetName() string {
	return "mockEmailService"
}

func (mes mockEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, mes.GetName(), "SendPlainTextEmail")
	defer span.End()
	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("TO:\t%v\n\n", to)

	fmt.Printf("SUBJECT:\t%s\n\n", subject)

	fmt.Printf("BODY:\t%s\n\n", body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}

type stackEmailService struct {
	messages []TestEmailMessage
}

func NewStackEmailService() *stackEmailService {
	messages := make([]TestEmailMessage, 0)
	return &stackEmailService{
		messages: messages,
	}
}

func (stackEmailService) GetName() string {
	return "stackEmailService"
}

func (ses *stackEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ses.GetName(), "SendPlainTextEmail")
	defer span.End()
	message := TestEmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	ses.messages = append(ses.messages, message)
	return nil
}

func (ses *stackEmailService) PopMessage() (TestEmailMessage, bool) {
	numMessages := len(ses.messages)
	if numMessages == 0 {
		return TestEmailMessage{}, false
	}
	message := ses.messages[numMessages-1]      // get the last message
	ses.messages = ses.messages[:numMessages-1] // save the array with the poped message clipped off
	return message, true
}
