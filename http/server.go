package http

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	corefactory "github.com/calvine/goauth/core/factory"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/factory"
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

	cachedJWTValidatorDuration time.Duration = time.Minute * 5 // TODO: make configurable?
)

type server struct {
	logger              *zap.Logger
	publicURL           string
	serviceName         string
	loginService        services.LoginService
	userService         services.UserService
	emailService        services.EmailService
	tokenService        services.TokenService
	appService          services.AppService
	jsmService          services.JWTSigningMaterialService
	jwtFactory          corefactory.JWTFactory
	jwtValidatorFactory corefactory.JWTValidatorFactory
	staticFS            *http.FileSystem
	templateFS          *embed.FS
	Mux                 *chi.Mux
	// allowedTokenSigningAlgorithmTypes tells us the jwt signing material to pull from the data store for token signing
	allowedTokenSigningAlgorithmTypes []jwt.JWTSigningAlgorithmFamily
	// validatorCache is a cache token validators so that we can quickly access them as needed
	validatorCache jwtValidatorCache
}

type HTTPServerOptions struct {
	Logger *zap.Logger
	// PublicURL is the URL used to create links to this service.
	// It needs to be the URL that can be used to make links for things like password reset,
	// and links for use in emails / notifications. For convention sake the url will NOT end with a slash.
	PublicURL                  string
	ServiceName                string
	LoginService               services.LoginService
	UserService                services.UserService
	EmailService               services.EmailService
	TokenService               services.TokenService
	AppService                 services.AppService
	JsmService                 services.JWTSigningMaterialService
	StaticFS                   *http.FileSystem
	TemplateFS                 *embed.FS
	TokenSigningAlgorithmTypes []jwt.JWTSigningAlgorithmFamily
	DefaultJWTValidatorOptions jwt.JWTValidatorOptions
}

func (hso HTTPServerOptions) Validate() errors.RichError {
	if hso.Logger == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.Logger", true)
	}
	if len(hso.PublicURL) == 0 {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.PublicURL", true)
	}
	if len(hso.ServiceName) == 0 {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.ServiceName", true)
	}
	if hso.LoginService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.LoginService", true)
	}
	if hso.UserService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.UserService", true)
	}
	if hso.EmailService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.EmailService", true)
	}
	if hso.TokenService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.TokenService", true)
	}
	if hso.AppService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.AppService", true)
	}
	if hso.JsmService == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.JsmService", true)
	}
	if hso.StaticFS == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.StaticFS", true)
	}
	if hso.TemplateFS == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.TemplateFS", true)
	}
	if hso.TokenSigningAlgorithmTypes == nil {
		return coreerrors.NewNilParameterNotAllowedError("http.NewServer", "options.TokenSigningAlgorithmTypes", true)
	}
	return nil
}

// TODO: implement a validate method for HTTPServerOptions

func NewServer(ctx context.Context, options HTTPServerOptions) (server, errors.RichError) {
	mux := chi.NewMux()
	err := options.Validate()
	if err != nil {
		if options.Logger != nil {
			options.Logger.Error("HttpServerOptions validation failed", zap.Reflect("err", err))
		}
		return server{}, err
	}
	jvf := factory.NewJWTValidatorFactory(options.DefaultJWTValidatorOptions)
	s := server{
		logger:                            options.Logger,
		serviceName:                       options.ServiceName,
		loginService:                      options.LoginService,
		userService:                       options.UserService,
		emailService:                      options.EmailService,
		tokenService:                      options.TokenService,
		appService:                        options.AppService,
		jsmService:                        options.JsmService,
		jwtValidatorFactory:               jvf,
		staticFS:                          options.StaticFS,
		templateFS:                        options.TemplateFS,
		Mux:                               mux,
		allowedTokenSigningAlgorithmTypes: options.TokenSigningAlgorithmTypes,
	}
	jwtSigningMaterial := make([]models.JWTSigningMaterial, 0, 3)
	for _, alg := range options.TokenSigningAlgorithmTypes {
		results, err := s.jsmService.GetValidJWTSigningMaterialByAlgorithmType(ctx, s.logger, alg, "server startup")
		if err != nil {
			s.logger.Error("call for jwt signing material failed!", zap.Any("error", err))
			return s, err
		}
		jwtSigningMaterial = append(jwtSigningMaterial, results...)
	}
	if len(jwtSigningMaterial) == 0 {
		fields := map[string]interface{}{
			"algorithmTypes": options.TokenSigningAlgorithmTypes,
		}
		err := coreerrors.NewNoJWTSigningMaterialFoundError(fields, true)
		return s, err
	}
	tokenSigners := make([]jwt.Signer, 0, 3)
	s.validatorCache = make(jwtValidatorCache)
	for _, jsm := range jwtSigningMaterial {
		options.Logger.Info("building signer for jwt signig material", zap.String("jsmID", jsm.ID), zap.String("jsmKeyID", jsm.KeyID))
		signer, err := jsm.ToSigner()
		if err != nil {
			options.Logger.Error("failed to create jwt signer from jwt signing material", zap.Any("error", err))
			return s, err
		}
		/* TODO: how do i select the right kind of key for signing? for now random select a key for signing?
		also any key id not present will be pulled at the time it is used, and if its not in the repo we throw an error...

		Change of plans. I was way overthinking this. later on we can have a list of alogrithms used and that will be fine...
		key lookups for verifying a token should be done via the jwt signing material service. I am going to leave all of this commentary here for now
		but eventually I will remove it once its all implemented completely.
		*/
		tokenSigners = append(tokenSigners, signer)
	}
	s.jwtFactory, err = factory.NewJWTFactory(options.ServiceName, tokenSigners)
	if err != nil {
		return s, err
	}
	return s, nil
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
