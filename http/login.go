package http

import (
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/viewmodels"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FIXME: add proper error handling like in register file
func (s *server) handleLoginGet() http.HandlerFunc {
	var (
		once        sync.Once
		templateErr error
	)
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()
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
		token, err := s.getNewSCRFToken(ctx, logger)
		if err != nil {
			errorMsg := "failed to create new CSRF token"
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		templateRenderError := loginPageTemplate.Execute(rw, viewmodels.LoginTemplateData{
			CSRF: token.Value,
		})
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
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()
		data := requestData{}
		data.CSRFToken = r.FormValue("csrf_token")
		data.Email = r.FormValue("email")
		data.Password = r.FormValue("password")

		_, err := s.retreiveCSRFToken(ctx, logger, data.CSRFToken)
		if err != nil {
			errorMsg := "failed to retreive CSRF token"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		_, err = s.loginService.LoginWithPrimaryContact(ctx, s.logger, data.Email, core.CONTACT_TYPE_EMAIL, data.Password, "login post handler")
		if err != nil {
			// TODO: add switch case here to set error message and re render login page for certain error messages
			errorMsg := "login attempt failed"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			http.Error(rw, err.GetErrorMessage(), http.StatusUnauthorized)
			return
		}
		// TODO: finish the login code...
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
