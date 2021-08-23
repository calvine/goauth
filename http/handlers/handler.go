package handlers

import "github.com/calvine/goauth/core/services"

type HttpHandler struct {
	loginService services.LoginService
	emailService services.EmailService
}
