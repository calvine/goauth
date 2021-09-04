package mux

import (
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/http/login"
	"github.com/calvine/goauth/http/passwordreset"
	"github.com/go-chi/chi/v5"
)

type httpMux struct {
	loginService services.LoginService
	emailService services.EmailService
	Mux          *chi.Mux
}

func NewHttpMux(loginService services.LoginService, emailService services.EmailService) httpMux {
	mux := chi.NewRouter()
	return httpMux{loginService, emailService, mux}
}

func (hh *httpMux) BuildRoutes() {
	hh.Mux.Route("/auth", func(r chi.Router) {
		r.Route("/login", func(r chi.Router) {
			// this is the route for the login page
			r.Get("/", login.GetLoginHandler)
			// this is the post target for the login page
			r.Post("/", login.PostLoginHandler)
		})
		r.Route("/resetpassword", func(r chi.Router) {
			r.Get("/{passwordResetToken}", passwordreset.GetPasswordResetHandler)
			r.Post("/reset", passwordreset.PostPasswordResetHandler)
		})
	})
}
