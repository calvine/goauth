package http

import (
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FIXME: add proper error handling like in register file
func (s *server) handleLoginGet() http.HandlerFunc {
	var (
		once        sync.Once
		templateErr error
	)
	type requestData struct {
		CSRFToken string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			if loginPageTemplate == nil {
				loginPageTemplate, templateErr = parseTemplateFromEmbedFS(loginPageTemplatePath, loginPageName, s.templateFS)
			}
		})
		if templateErr != nil {
			http.Error(rw, templateErr.Error(), http.StatusInternalServerError)
			return
		}
		// TODO: make CSRF token life span configurable
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		token, err := models.NewToken("", models.TokenTypeCSRF, constants.Default_CSRF_Token_Duration)
		if err != nil {
			span.RecordError(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		err = s.tokenService.PutToken(ctx, logger, token)
		if err != nil {
			span.RecordError(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		templateRenderError := loginPageTemplate.Execute(rw, requestData{token.Value})
		if templateRenderError != nil {
			span.RecordError(err)
			err = coreerrors.NewFailedTemplateRenderError(loginPageName, templateRenderError, true)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleLoginPost() http.HandlerFunc {
	// var (
	// 	once          sync.Once
	// 	loginTemplate *template.Template
	// 	templateErr   error
	// 	templatePath  string = "http/templates/login.tmpl"
	// )
	type requestData struct {
		CSRFToken string
		Email     string
		Password  string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		data := requestData{}
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		data.CSRFToken = r.FormValue("csrf_token")
		data.Email = r.FormValue("email")
		data.Password = r.FormValue("password")

		_, err := s.tokenService.GetToken(ctx, logger, data.CSRFToken, models.TokenTypeCSRF)
		if err != nil {
			http.Error(rw, err.GetErrorMessage(), http.StatusBadRequest)
			return
		}
		err = s.tokenService.DeleteToken(ctx, logger, data.CSRFToken)
		if err != nil {
			logger.Warn("unable to delete CSRF token from data store", zap.Reflect("error", err))
			// uh of the token was not deleted! need to log this...
		}
		_, err = s.loginService.LoginWithPrimaryContact(ctx, s.logger, data.Email, core.CONTACT_TYPE_EMAIL, data.Password, "login post handler")
		if err != nil {
			http.Error(rw, err.GetErrorMessage(), http.StatusUnauthorized)
			return
		}
		rw.Header().Add("test-header", "JWT")
		http.Redirect(rw, r, "/static/hooray.html", http.StatusFound)
	}
}

// func (s *server) handleAuthGet() http.HandlerFunc {
// 	var (
// 		once          sync.Once
// 		loginTemplate *template.Template
// 		templateErr   error
// 		templatePath  string = "http/templates/login.tmpl"
// 	)
// 	type requestData struct {
// 		clientID      string
// 		codeChallenge string // PKCE
// 		redirectURI   string
// 		responseType  string
// 		scope         string
// 		state         string
// 	}
// 	return func(rw http.ResponseWriter, r *http.Request) {
// 		once.Do(func() {
// 			templateFileData, err := s.templateFS.ReadFile(templatePath)
// 			templateErr = err
// 			if templateErr == nil {
// 				loginTemplate, templateErr = template.New("loginPage").Parse(string(templateFileData))
// 			}
// 		})
// 		cookie, err := r.Cookie(loginCookieName)
// 		if err == http.ErrNoCookie {
// 			// handle not logged in
// 		}

// 		http.SetCookie(rw, &http.Cookie{
// 			Name:     loginCookieName,
// 			Value:    "session token here",
// 			Expires:  time.Now().Add(time.Hour * 24 * 7),
// 			SameSite: http.SameSiteLaxMode,
// 			Secure:   true,
// 			HttpOnly: true,
// 		})

// 	}
// }
