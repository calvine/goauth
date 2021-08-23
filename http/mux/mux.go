package mux

import (
	"time"

	"github.com/calvine/goauth/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type httpMux struct {
	httpService handlers.HttpHandler
	Mux         *chi.Mux
}

func NewHttpMux(httpHandler handlers.HttpHandler) httpMux {
	mux := chi.NewRouter()
	return httpMux{httpHandler, mux}
}

func (hh *httpMux) BuildRoutes() {
	hh.Mux.Use(
		middleware.Recoverer,
		middleware.Timeout(time.Second*5),
		middleware.RequestID,
		middleware.RealIP,
	)
	hh.Mux.Route("/auth", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Route("/login", func(r chi.Router) {
			// this is the route for the login page
			r.Get("/", hh.httpService.GetLoginHandler)
			// this is the post target for the login page
			r.Post("/", hh.httpService.PostLoginHandler)
		})
		r.Route("/resetpassword", func(r chi.Router) {
			// this is the route for the password reset page
			r.Get("/{passwordResetToken}", hh.httpService.GetPasswordResetHandler)
			// this is the post endpoint for the password reset page
			r.Post("/{passwordResetToken}", hh.httpService.PostPasswordResetHandler)
		})
	})
}
