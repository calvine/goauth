package service

import (
	"context"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	coreServices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/internal/constants"
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
	From    string
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
		return newStackEmailService(), nil
	case SMTPEmailService: // TODO: implement this...
		castOptions, ok := options.(SMTPEmailServiceOptions)
		if !ok {
			return nil, coreerrors.NewInvalidSMTPEmailOptionsError("failed to cast options to SMTPEmailServiceOptions", nil, true)
		}
		return newSMTPEmailService(castOptions)
	default:
		return nil, coreerrors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}

type noopEmailService struct{}

func (noopEmailService) GetName() string {
	return "noopEmailService"
}

func (ne noopEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, from, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ne.GetName(), "SendPlainTextEmail")
	defer span.End()
	return nil
}

type mockEmailService struct{}

func (mockEmailService) GetName() string {
	return "mockEmailService"
}

func (mes mockEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, from, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, mes.GetName(), "SendPlainTextEmail")
	defer span.End()

	if from == "" {
		logger.Info("no from address supplied, so using default")
		from = constants.NoReplyEmailAddress
	}

	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("FROM:\t%v\n\n", from)

	fmt.Printf("TO:\t%v\n\n", to)

	fmt.Printf("SUBJECT:\t%s\n\n", subject)

	fmt.Printf("BODY:\t%s\n\n", body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}

type stackEmailService struct {
	messages []TestEmailMessage
}

func newStackEmailService() *stackEmailService {
	messages := make([]TestEmailMessage, 0)
	return &stackEmailService{
		messages: messages,
	}
}

func (stackEmailService) GetName() string {
	return "stackEmailService"
}

func (ses *stackEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, from, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ses.GetName(), "SendPlainTextEmail")
	defer span.End()

	if from == "" {
		logger.Info("no from address supplied, so using default")
		from = constants.NoReplyEmailAddress
	}

	message := TestEmailMessage{
		From:    from,
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

type smtpEmailService struct {
	host     string
	port     string
	user     string
	password string
}

type SMTPEmailServiceOptions struct {
	Host     string
	Port     string
	User     string
	Password string
}

func newSMTPEmailService(options SMTPEmailServiceOptions) (*smtpEmailService, errors.RichError) {
	// validate host
	if options.Host == "" {
		return nil, coreerrors.NewInvalidSMTPEmailOptionsError("host must have a value", nil, true)
	}

	// validate port
	numPort, err := strconv.Atoi(options.Port)
	if err != nil {
		fields := map[string]interface{}{
			"port": options.Port,
		}
		return nil, coreerrors.NewInvalidSMTPEmailOptionsError("port must be a numeric value", fields, true)
	}
	if numPort <= 0 {
		fields := map[string]interface{}{
			"port": options.Port,
		}
		return nil, coreerrors.NewInvalidSMTPEmailOptionsError("port must be a value greater than zero", fields, true)
	}

	return &smtpEmailService{
		host:     options.Host,
		port:     options.Port,
		user:     options.User,
		password: options.Password,
	}, nil
}

func (smtpEmailService) GetName() string {
	return "smtpEmailService"
}

func (ses *smtpEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, from, subject, body string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ses.GetName(), "SendPlainTextEmail")
	defer span.End()

	if from == "" {
		logger.Info("no from address supplied, so using default")
		from = constants.NoReplyEmailAddress
	}

	addr := fmt.Sprintf("%s:%s", ses.host, ses.port)

	auth := smtp.PlainAuth("", from, ses.password, ses.host)

	err := smtp.SendMail(addr, auth, from, to, []byte(body))

	if err != nil {
		rErr := coreerrors.NewFailedToSendEmailError(err, nil, true)
		logger.Error("failed to send email", zap.Reflect("error", rErr))
		return rErr
	}

	return nil
}
