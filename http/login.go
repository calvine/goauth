package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
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
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
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
		callback := r.URL.Query().Get("cb")

		_, err := s.retreiveCSRFToken(ctx, logger, csrfToken)
		if err != nil {
			errorMsg := "failed to retreive CSRF token"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		_, err = s.loginService.LoginWithPrimaryContact(ctx, s.logger, email, core.CONTACT_TYPE_EMAIL, password, "login post handler")
		if err != nil {
			var errorMsg string
			switch err.GetErrorCode() {
			case coreerrors.ErrCodeLoginContactNotPrimary:
				errorMsg = "no account found or incorrect password"
			case coreerrors.ErrCodeLoginFailedWrongPassword:
				errorMsg = "no account found or incorrect password"
			case coreerrors.ErrCodeLoginPrimaryContactNotConfirmed:
				errorMsg = "contact not confirmed"
				// TODO: Should we show a send confirmation link when this is shown?
			case coreerrors.ErrCodeNoUserFound:
				errorMsg = "no account found or incorrect password"
			default:
				errorMsg = "login attempt failed"
				logger.Error(errorMsg, zap.Reflect("error", err))
				apptelemetry.SetSpanError(&span, err, errorMsg)
				// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
				redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
				return
			}
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			newCSRFToken, err := s.getNewSCRFToken(ctx, logger)
			if err != nil {
				errorMsg := "failed to create new CSRF token"
				apptelemetry.SetSpanError(&span, err, errorMsg)
				redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
				// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
				return
			}
			templateRenderError := loginPageTemplate.Execute(rw, viewmodels.LoginTemplateData{
				CSRFToken: newCSRFToken.Value,
				Email:     email,
				ErrorMsg:  errorMsg,
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
		if callback != "" {
			http.Redirect(rw, r, callback, http.StatusFound)
		} else {
			http.Redirect(rw, r, "/static/hooray.html", http.StatusFound)
		}
	}
}

func (s *server) handleMagicLoginGet() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()

		q := r.URL.Query()
		magicLoginToken := q.Get("m")
		if magicLoginToken == "" {
			// return bad request?
			errorMsg := "no magic token provided"
			err := coreerrors.NewNoMagicLoginTokenFoundError(true)
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, errorMsg, http.StatusBadRequest)
			return
		}
		token, err := s.tokenService.GetToken(ctx, logger, magicLoginToken, models.TokenTypeMagicLogin)
		if err != nil {
			// return error based on error returned
			var errorMsg string
			switch err.GetErrorCode() {
			default:
				errorMsg = "an unexpected error occurred"
			}
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, errorMsg, http.StatusBadRequest)
			return
		}
		userID := token.TargetID
		if userID == "" {
			// return bad request or not found?
			errorMsg := "magic token not associated with a user"
			err := coreerrors.NewMagicLoginTokenNoUserIDError(magicLoginToken, true)
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, "token provided was not valid", http.StatusBadRequest)
			return
		}
		// login successfull?
		_, err = s.userService.GetUser(ctx, logger, userID, "magic login")
		if err != nil {
			// return proper error based on error returned above
		}
		// TODO: create JWT and store cookie
	}
}

type setLoginStateOptions struct {
	user     models.User
	signer   jwt.Signer
	duration time.Duration
}

func (server) setLoginState(ctx context.Context, logger *zap.Logger, rw http.ResponseWriter, options setLoginStateOptions) errors.RichError {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	jwToken := jwt.NewUnsignedJWT(jwt.HS256, "goauth", []string{}, options.user.ID, options.duration, time.Now())
	span.AddEvent("unsigned jwt created")
	// hmacOptions, err := jwt.NewHMACSigningOptions("test")
	// if err != nil {
	// 	errorMsg := "building signing options failed"
	// 	logger.Error(errorMsg, zap.Reflect("error", err))
	// 	apptelemetry.SetSpanOriginalError(&span, err, errorMsg)
	// 	return err
	// }
	encodedJWT, err := jwToken.SignAndEncode(options.signer)
	if err != nil {
		errorMsg := "signing and encoding jwt failed"
		logger.Error(errorMsg, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, errorMsg)
		return err
	}
	span.AddEvent("JWT signed and encoded")
	// encodedHWT, err := jwt.SignAndEncode()
	cookie := http.Cookie{
		Name:     constants.LoginCookieName,
		Value:    encodedJWT,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // TODO: Evaluate what the best setting is for this. Starting with strict.
	}
	http.SetCookie(rw, &cookie)
	span.AddEvent("sign in cookie set")
	return nil
}

func (server) getAuthStatus(ctx context.Context, logger *zap.Logger, r *http.Request) (core.AuthStatus, errors.RichError) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	authCookie, err := r.Cookie(constants.LoginCookieName)
	if err != nil {
		// if err == http.ErrNoCookie {
		// according to the r.Cookie method the only possible error is ErrNoCookie,
		// so any error here will be treated as if the cookie was just not there.
		logger.Debug("request did not contain auth cookie")
		return core.Unauthenticated, nil
	}
	jwtParts, rErr := jwt.SplitEncodedJWT(authCookie.Value)
	if rErr != nil {
		errMsg := "failed to split encoded jwt from auth cookie!"
		logger.Error(errMsg, zap.Reflect("err", rErr), zap.String("cookie_value", authCookie.Value))
		apptelemetry.SetSpanOriginalError(&span, rErr, errMsg)
		return core.Invalid, rErr
	}
	header, rErr := jwt.DecodeHeader(jwtParts[0])
	if rErr != nil {
		errMsg := "failed to decode header from auth cookie jwt"
		logger.Error(errMsg, zap.Reflect("err", rErr))
		apptelemetry.SetSpanOriginalError(&span, rErr, errMsg)
	}
	// get / make jwt validator...
	// if valid the we are authenticated
	// handle expired?
	if len(header.KeyID) == 0 {

	} else {

	}
	return core.Unauthenticated, nil
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
