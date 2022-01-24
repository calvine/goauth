package http

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/services"
	mymiddleware "github.com/calvine/goauth/http/middleware"
	"github.com/calvine/richerror/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

var (
	accountRegisteredPageTemplate *template.Template
	authPageTemplate              *template.Template
	errorPageTemplate             *template.Template
	loginPageTemplate             *template.Template
	redirectPageTemplate          *template.Template
	registerPageTemplate          *template.Template
)

const (
	accountRegisteredPageName         = "accountRegistered"
	accountRegisteredPageTemplatePath = "http/templates/accountregistered.html.tmpl"
	authPageName                      = "auth"
	authPageTemplatePath              = "http/templates/auth.html.tmpl"
	errorPageName                     = "error"
	errorPageTemplatePath             = "http/templates/error.html.tmpl"
	loginPageName                     = "login"
	loginPageTemplatePath             = "http/templates/login.html.tmpl"
	redirectPageName                  = "redirect"
	redirectPageTemplatePath          = "http/templates/redirect.html.tmpl"
	registerPageName                  = "register"
	registerPageTemplatePath          = "http/templates/register.html.tmpl"
)

type server struct {
	logger       *zap.Logger
	loginService services.LoginService
	userService  services.UserService
	emailService services.EmailService
	tokenService services.TokenService
	appService   services.AppService
	jsmService   services.JWTSigningMaterialService
	staticFS     *http.FileSystem
	templateFS   *embed.FS
	Mux          *chi.Mux
}

type HTTPServerOptions struct {
	logger       *zap.Logger
	loginService services.LoginService
	userService  services.UserService
	emailService services.EmailService
	tokenService services.TokenService
	appService   services.AppService
	jsmService   services.JWTSigningMaterialService
	staticFS     *http.FileSystem
	templateFS   *embed.FS
	Mux          *chi.Mux
}

func NewServer(logger *zap.Logger, loginService services.LoginService, userService services.UserService, emailService services.EmailService, tokenService services.TokenService, appService services.AppService, jsms services.JWTSigningMaterialService, staticFS *http.FileSystem, templateFS *embed.FS) server {
	mux := chi.NewRouter()
	return server{logger, loginService, userService, emailService, tokenService, appService, jsms, staticFS, templateFS, mux}
}

func (hh *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.Mux.ServeHTTP(w, r)
}

func (hh *server) ParseTemplates() errors.RichError {
	var rErr errors.RichError
	// parse accountregistered page template
	accountRegisteredPageTemplate, rErr = parseTemplateFromEmbedFS(accountRegisteredPageTemplatePath, accountRegisteredPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	// parse auth page template
	authPageTemplate, rErr = parseTemplateFromEmbedFS(authPageTemplatePath, authPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	// parse error page template
	errorPageTemplate, rErr = parseTemplateFromEmbedFS(errorPageTemplatePath, errorPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	// parse login page template
	loginPageTemplate, rErr = parseTemplateFromEmbedFS(loginPageTemplatePath, loginPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	// parse redirect page template
	redirectPageTemplate, rErr = parseTemplateFromEmbedFS(redirectPageTemplatePath, redirectPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	// parse register page template
	registerPageTemplate, rErr = parseTemplateFromEmbedFS(registerPageTemplatePath, registerPageName, hh.templateFS)
	if rErr != nil {
		return rErr
	}
	return nil
}

func (hh *server) BuildRoutes() {
	hh.Mux.Use(
		// middleware.Recoverer,
		mymiddleware.InitializeRequest(hh.logger),
		middleware.Timeout(time.Second*60),
		middleware.RealIP,
	)
	hh.Mux.Route("/error", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Get("/", otelhttp.NewHandler(hh.handleErrorGet(), "GET /error").ServeHTTP)
	})
	hh.Mux.Route("/auth", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Route("/login", func(r chi.Router) {
			// this is the route for the login page
			r.Get("/", otelhttp.NewHandler(hh.handleLoginGet(), "GET /auth/login").ServeHTTP) //addTrace(hh.handleLoginGet(), "GET /auth/login"))
			// this is the post target for the login page
			r.Post("/", otelhttp.NewHandler(hh.handleLoginPost(), "POST /auth/login").ServeHTTP)
		})
		r.Route("/magiclogin", func(r chi.Router) {
			// this is the route for magic login
			r.Get("/", otelhttp.NewHandler(hh.handleMagicLoginGet(), "GET /auth/magiclogin").ServeHTTP)
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

func parseTemplateFromEmbedFS(path string, name string, fs *embed.FS) (*template.Template, errors.RichError) {
	templateFileData, err := fs.ReadFile(path)
	if err != nil {
		rErr := coreerrors.NewFailedTemplateParseTemplatNotFoundError(name, path, err, true)
		return nil, rErr
	}
	return parseTemplate(name, string(templateFileData))
}

func parseTemplate(name string, templateString string) (*template.Template, errors.RichError) {
	template, err := template.New(name).Parse(templateString)
	if err != nil {
		rErr := coreerrors.NewFailedTemplateParseError(name, err, true)
		return nil, rErr
	}
	return template, nil
}

func renderTemplate(rw http.ResponseWriter, template *template.Template, data interface{}) errors.RichError {
	err := template.Execute(rw, data)
	if err != nil {
		rErr := coreerrors.NewFailedTemplateRenderError(template.Name(), err, true)
		return rErr
	}
	return nil
}
