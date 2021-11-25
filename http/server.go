package http

import (
	"embed"
	"net/http"
	"time"

	"github.com/calvine/goauth/core/services"
	mymiddleware "github.com/calvine/goauth/http/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type server struct {
	logger       *zap.Logger
	loginService services.LoginService
	userService  services.UserService
	emailService services.EmailService
	tokenService services.TokenService
	staticFS     *http.FileSystem
	templateFS   *embed.FS
	Mux          *chi.Mux
}

func NewServer(logger *zap.Logger, loginService services.LoginService, userService services.UserService, emailService services.EmailService, tokenService services.TokenService, staticFS *http.FileSystem, templateFS *embed.FS) server {
	mux := chi.NewRouter()
	return server{logger, loginService, userService, emailService, tokenService, staticFS, templateFS, mux}
}

func (hh *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.Mux.ServeHTTP(w, r)
}

func (hh *server) BuildRoutes() {
	hh.Mux.Use(
		// middleware.Recoverer,
		mymiddleware.InitializeRequest(hh.logger),
		middleware.Timeout(time.Second*60),
		middleware.RealIP,
	)
	hh.Mux.Route("/auth", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Route("/login", func(r chi.Router) {
			// this is the route for the login page
			r.Get("/", otelhttp.NewHandler(hh.handleLoginGet(), "GET /auth/login").ServeHTTP) //addTrace(hh.handleLoginGet(), "GET /auth/login"))
			// this is the post target for the login page
			r.Post("/", otelhttp.NewHandler(hh.handleLoginPost(), "POST /auth/login").ServeHTTP)
		})
		r.Route("/resetpassword", func(r chi.Router) {
			// this is the route for the password reset page
			r.Get("/{passwordResetToken}", otelhttp.NewHandler(hh.handlePasswordResetGet(), "GET /resetpassword/{passwordResetToken}").ServeHTTP)
			// this is the post endpoint for the password reset page
			r.Post("/submitpasswordreset", otelhttp.NewHandler(hh.handlePasswordResetPost(), "POST /resetpassword/submitpasswordreset").ServeHTTP)
		})
	})
	hh.Mux.Route("/user", func(r chi.Router) {
		r.Get("/register", otelhttp.NewHandler(hh.handleRegisterGet(), "GET /user/register").ServeHTTP)
		r.Post("/register", otelhttp.NewHandler(hh.handleRegisterPost(), "POST /user/register").ServeHTTP)

		r.Get("/confirmcontact/{confirmationToken}", otelhttp.NewHandler(hh.handleConfirmContactGet(), "GET /user/confirmcontact/{confirmationToken}").ServeHTTP)
	})
	hh.Mux.Route("/app", func(r chi.Router) {
		r.Get("/manage", otelhttp.NewHandler(nil, "GET /app/manage").ServeHTTP)

		r.Get("/manage/{appID}", otelhttp.NewHandler(nil, "GET /app/manage/{appID}").ServeHTTP)
		r.Post("/manage/{appID}", otelhttp.NewHandler(nil, "POST /app/manage/{appID}").ServeHTTP)
	})
	fs := http.FileServer(*hh.staticFS)
	// static files
	hh.Mux.Handle("/static/*", fs)
}
