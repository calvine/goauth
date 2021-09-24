package service

import (
	"context"
	"fmt"

	coreerrors "github.com/calvine/goauth/core/errors"
	coreServices "github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

const (
	MockEmailService = "mock"
	NoOpEmailService = "noop"
	SMTPEmailService = "smtp"
)

func NewEmailService(serviceType string, options interface{}) (coreServices.EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	case NoOpEmailService:
		return noopEmailService{}, nil
	// case SMTPEmailService: // TODO: implement this...
	default:
		return nil, coreerrors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}

type noopEmailService struct{}

func (noopEmailService) GetName() string {
	return "noopEmailService"
}

func (ne noopEmailService) SendPlainTextEmail(ctx context.Context, to []string, subject, body string) errors.RichError {
	return nil
}

type mockEmailService struct{}

func (mockEmailService) GetName() string {
	return "mockEmailService"
}

func (mse mockEmailService) SendPlainTextEmail(ctx context.Context, to []string, subject, body string) errors.RichError {
	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("TO:\t%v\n\n", to)

	fmt.Printf("SUBJECT:\t%s\n\n", subject)

	fmt.Printf("BODY:\t%s\n\n", body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}
