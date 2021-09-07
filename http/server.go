package http

import (
	"embed"
	"net/http"
	"time"

	"github.com/calvine/goauth/core/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type server struct {
	loginService services.LoginService
	emailService services.EmailService
	staticFS     *http.FileSystem
	templateFS   *embed.FS
	Mux          *chi.Mux
}

func NewServer(loginService services.LoginService, emailService services.EmailService, staticFS *http.FileSystem, templateFS *embed.FS) server {
	mux := chi.NewRouter()
	return server{loginService, emailService, staticFS, templateFS, mux}
}

func (hh *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.Mux.ServeHTTP(w, r)
}

func (hh *server) BuildRoutes() {
	hh.Mux.Use(
		// middleware.Recoverer,
		middleware.Timeout(time.Second*5),
		middleware.RequestID,
		middleware.RealIP,
	)
	hh.Mux.Route("/auth", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Route("/login", func(r chi.Router) {
			// this is the route for the login page
			r.Get("/", hh.handleLoginGet())
			// this is the post target for the login page
			r.Post("/", hh.handleLoginPost())
		})
		r.Route("/resetpassword", func(r chi.Router) {
			// this is the route for the password reset page
			r.Get("/{passwordResetToken}", hh.handlePasswordResetGet())
			// this is the post endpoint for the password reset page
			r.Post("/{passwordResetToken}", hh.handlePasswordResetPost())
		})
	})
	fs := http.FileServer(*hh.staticFS)
	// static files
	hh.Mux.Handle("/static/*", fs)
}
