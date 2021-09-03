package http

import (
	"github.com/calvine/goauth/core/services"
)

type httpHandler struct {
	loginService services.LoginService
	emailService services.EmailService
}

func NewHttpHandler(loginService services.LoginService, emailService services.EmailService) httpHandler {
	return httpHandler{loginService, emailService}
}

func (hh *httpHandler) Init() {

}
