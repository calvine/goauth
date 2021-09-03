package service

import (
	"fmt"

	coreerrors "github.com/calvine/goauth/core/errors"
	coreServices "github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

const (
	MockEmailService = "mock"
	SMTPEmailService = "smtp"
)

func NewEmailService(serviceType string, options interface{}) (coreServices.EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	// case SMTPEmailService: // TODO: implement this...
	default:
		return nil, coreerrors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}

type mockEmailService struct{}

func (mse mockEmailService) SendPlainTextEmail(to []string, subject, body string) errors.RichError {
	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("TO:\t%v\n\n", to)

	fmt.Printf("SUBJECT:\t%s\n\n", subject)

	fmt.Printf("BODY:\t%s\n\n", body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}
