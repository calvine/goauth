package http

import (
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
	"github.com/calvine/goauth/http/internal/viewmodels"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// FIXME: add proper error handling like in register file
func (s *server) handleLoginGet() http.HandlerFunc {
	var (
		once        sync.Once
		templateErr errors.RichError
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
			errorMsg := "initial parsing of template failed"
			logger.Error(errorMsg, zap.Reflect("error", templateErr))
			apptelemetry.SetSpanError(&span, templateErr, errorMsg)
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
			CSRFToken: token.Value,
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
	var (
		once        sync.Once
		templateErr errors.RichError
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
			errorMsg := "initial parsing of template failed"
			logger.Error(errorMsg, zap.Reflect("error", templateErr))
			apptelemetry.SetSpanError(&span, templateErr, errorMsg)
			http.Error(rw, templateErr.Error(), http.StatusInternalServerError)
			return
		}
		csrfToken := r.FormValue("csrf_token")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_, err := s.retreiveCSRFToken(ctx, logger, csrfToken)
		if err != nil {
			errorMsg := "failed to retreive CSRF token"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		_, err = s.loginService.LoginWithPrimaryContact(ctx, s.logger, email, core.CONTACT_TYPE_EMAIL, password, "login post handler")
		if err != nil {
			var errorMsg string
			switch err.GetErrorCode() {
			case coreerrors.ErrCodeLoginContactNotPrimary:
			case coreerrors.ErrCodeLoginFailedWrongPassword:
			case coreerrors.ErrCodeLoginPrimaryContactNotConfirmed:
			case coreerrors.ErrCodeNoUserFound:
			// TODO: add switch case here to set error message and re render login page for certain error messages
			default:
				errorMsg = "login attempt failed"
				logger.Error(errorMsg, zap.Reflect("error", err))
				apptelemetry.SetSpanError(&span, err, errorMsg)
				s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
				return
			}
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			newCSRFToken, err := s.getNewSCRFToken(ctx, logger)
			if err != nil {
				errorMsg := "failed to create new CSRF token"
				apptelemetry.SetSpanError(&span, err, errorMsg)
				s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
				return
			}
			templateRenderError := loginPageTemplate.Execute(rw, viewmodels.LoginTemplateData{
				CSRFToken:       newCSRFToken.Value,
				Email:           email,
				ErrorMsg:        errorMsg,
				HasErrorMessage: true,
				// Do Not Set Password!!!
			})
			if templateRenderError != nil {
				span.RecordError(err)
				err = coreerrors.NewFailedTemplateRenderError(loginPageName, templateRenderError, true)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		// TODO: finish the login code...
		authCookie := http.Cookie{
			Name: constants.LoginCookieName,
			// TODO: make a lightweight jwt for this
			Value:    "make a lightweight JWT for this I do not want to store a session id...",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(rw, &authCookie)
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
