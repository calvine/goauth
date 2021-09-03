package http

import (
	"github.com/calvine/goauth/core/services"
	"github.com/go-chi/chi/v5"
)

type httpHandler struct {
	loginService services.LoginService
	emailService services.EmailService
	mux          *chi.Mux
}

func NewHttpHandler(loginService services.LoginService, emailService services.EmailService) httpHandler {
	mux := chi.NewRouter()
	return httpHandler{loginService, emailService, mux}
}

func (hh *httpHandler) Init() {

}
