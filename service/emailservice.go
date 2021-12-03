package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models/email"
	coreServices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/internal/constants"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

const (
	MockEmailService  = "mock"
	NoOpEmailService  = "noop"
	StackEmailService = "stack"
	FSEmailService    = "fs"
	SMTPEmailService  = "smtp"
)

func NewEmailService(serviceType string, options interface{}) (coreServices.EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	case NoOpEmailService:
		return noopEmailService{}, nil
	case StackEmailService:
		return newStackEmailService(), nil
	case FSEmailService:
		castOptions, ok := options.(FSEmailServiceOptions)
		if !ok {
			return nil, coreerrors.NewInvalidSMTPEmailOptionsError("failed to cast options to FSEmailServiceOptions", nil, true)
		}
		return newFSEmailService(castOptions), nil
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

func (ne noopEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, message email.EmailMessage) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ne.GetName(), "SendPlainTextEmail")
	defer span.End()
	return nil
}

type mockEmailService struct{}

func (mockEmailService) GetName() string {
	return "mockEmailService"
}

func (mes mockEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, message email.EmailMessage) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, mes.GetName(), "SendPlainTextEmail")
	defer span.End()

	if message.From == "" {
		logger.Info("no from address supplied, so using default")
		message.From = constants.NoReplyEmailAddress
	}

	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("FROM:\t%v\n\n", message.From)

	fmt.Printf("TO:\t%v\n\n", message.To)

	fmt.Printf("SUBJECT:\t%s\n\n", message.Subject)

	fmt.Printf("BODY:\t%s\n\n", message.Body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}

type stackEmailService struct {
	messages []email.EmailMessage
}

func newStackEmailService() *stackEmailService {
	messages := make([]email.EmailMessage, 0)
	return &stackEmailService{
		messages: messages,
	}
}

func (stackEmailService) GetName() string {
	return "stackEmailService"
}

func (ses *stackEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, message email.EmailMessage) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ses.GetName(), "SendPlainTextEmail")
	defer span.End()

	if message.From == "" {
		logger.Info("no from address supplied, so using default")
		message.From = constants.NoReplyEmailAddress
	}

	ses.messages = append(ses.messages, message)
	return nil
}

func (ses *stackEmailService) PopMessage() (email.EmailMessage, bool) {
	numMessages := len(ses.messages)
	if numMessages == 0 {
		return email.EmailMessage{}, false
	}
	message := ses.messages[numMessages-1]      // get the last message
	ses.messages = ses.messages[:numMessages-1] // save the array with the poped message clipped off
	return message, true
}

type fsEmailService struct {
	messageDir string
}

type FSEmailServiceOptions struct {
	MessageDir string
}

func newFSEmailService(options FSEmailServiceOptions) *fsEmailService {
	err := os.MkdirAll(options.MessageDir, 0766)
	if err != nil {
		panic(fmt.Sprintf("debug fs email service failed to create path for files: %s", err.Error()))
	}
	return &fsEmailService{
		messageDir: options.MessageDir,
	}
}

func (fsEmailService) GetName() string {
	return "fsEmailService"
}

func (ses *fsEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, message email.EmailMessage) errors.RichError {
	messageFileName := fmt.Sprintf("%d_%s.json", time.Now().Unix(), message.Subject)
	filePath := path.Join(ses.messageDir, messageFileName)
	fileData, err := json.MarshalIndent(message, "", "    ")
	if err != nil {
		rErr := coreerrors.NewFailedToSendEmailError(err, nil, true)
		return rErr
	}
	err = ioutil.WriteFile(filePath, fileData, 0644)
	if err != nil {
		rErr := coreerrors.NewFailedToSendEmailError(err, nil, true)
		return rErr
	}
	return nil
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

func (ses *smtpEmailService) SendPlainTextEmail(ctx context.Context, logger *zap.Logger, message email.EmailMessage) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ses.GetName(), "SendPlainTextEmail")
	defer span.End()
	from := message.From
	if from == "" {
		logger.Info("no from address supplied, so using default")
		// TODO: do we want to remove the default and require a from address each time this method is called?
		from = constants.NoReplyEmailAddress
	}

	addr := fmt.Sprintf("%s:%s", ses.host, ses.port)

	auth := smtp.PlainAuth("", from, ses.password, ses.host)

	err := smtp.SendMail(addr, auth, from, message.To, []byte(message.Body))

	if err != nil {
		rErr := coreerrors.NewFailedToSendEmailError(err, nil, true)
		logger.Error("failed to send email", zap.Reflect("error", rErr))
		return rErr
	}

	return nil
}
