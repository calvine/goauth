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
	"github.com/calvine/goauth/core/normalization"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
	"github.com/calvine/goauth/http/internal/viewmodels"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//
const defaultJWTDuration time.Duration = time.Hour * 24 * 7 // 1 week as default // FIXME: make this configurable...

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
		rememberMeString := r.FormValue("remember_me")
		password := r.FormValue("password")
		callback := r.URL.Query().Get("cb")

		rememberMe, _ := normalization.ReadBoolValue(rememberMeString, true)

		_, err := s.retreiveCSRFToken(ctx, logger, csrfToken)
		if err != nil {
			errorMsg := "failed to retreive CSRF token"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		user, err := s.loginService.LoginWithPrimaryContact(ctx, s.logger, email, core.CONTACT_TYPE_EMAIL, password, "login post handler")
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
				CSRFToken:  newCSRFToken.Value,
				Email:      email,
				ErrorMsg:   errorMsg,
				RememberMe: rememberMe,
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
		setLoginStateOptions := setLoginStateOptions{
			user:       user,
			rememberMe: rememberMe,
		}
		_, _, err = s.setLoginState(ctx, logger, rw, setLoginStateOptions)
		if err != nil {
			errorMsg := "failed to set login state"
			logger.Error(errorMsg, zap.Reflect("err", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
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
		user, err := s.userService.GetUser(ctx, logger, userID, "magic login")
		if err != nil {
			// return proper error based on error returned above
			errorMsg := "failed to retreive user"
			logger.Error(errorMsg, zap.Reflect("err", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		loginAuthDuration := nullable.NullableDuration{HasValue: false, Value: 0}
		_, _, err = s.setLoginState(ctx, logger, rw, setLoginStateOptions{user, loginAuthDuration, false})
		if err != nil {
			errorMsg := "failed to set login state"
			logger.Error(errorMsg, zap.Reflect("err", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		// TODO: redirect?
	}
}

type setLoginStateOptions struct {
	user       models.User
	duration   nullable.NullableDuration
	rememberMe bool
}

// func (server) createJWT(sub string, validDuration time.Duration, signer jwt.Signer) (jwt.JWT, errors.RichError) {

// }

// Set login state sets a cookie on the current response that contains the encoded jwt as its valie. it also returns the JWT struct and the encoded JWT
func (s server) setLoginState(ctx context.Context, logger *zap.Logger, rw http.ResponseWriter, options setLoginStateOptions) (jwt.JWT, string, errors.RichError) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	// expireDuration should always have a value!
	cookieExpires := time.Time{}                                 // default there is no expiration for site login cookie.
	expireDuration := nullable.NullableDuration{HasValue: false} // default there is no expiration for site login cookie.
	if !options.rememberMe {                                     // is options.rememberMe is true then we go with default no expiration
		if options.duration.HasValue {
			expireDuration = options.duration
			cookieExpires = time.Now().Add(options.duration.Value)
		} else {
			expireDuration = nullable.NullableDuration{HasValue: true, Value: defaultJWTDuration} // setting a default incase none was provided...
			cookieExpires = time.Now().Add(defaultJWTDuration)
		}
	}
	jwToken, encodedJWT, err := s.jwtFactory.NewSignedJWT(options.user.ID, []string{s.serviceName}, expireDuration)
	if err != nil {
		errorMsg := "signing and encoding jwt failed"
		logger.Error(errorMsg, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, errorMsg)
		return jwt.JWT{}, "", err
	}
	span.AddEvent("JWT signed and encoded")
	// encodedHWT, err := jwt.SignAndEncode()
	cookie := http.Cookie{
		Expires:  cookieExpires,
		Name:     constants.LoginCookieName,
		Value:    encodedJWT,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // TODO: Evaluate what the best setting is for this. Starting with strict.
	}
	http.SetCookie(rw, &cookie)
	span.AddEvent("sign in cookie set")
	return jwToken, encodedJWT, nil
}

func (s server) getAuthStatus(ctx context.Context, logger *zap.Logger, r *http.Request, rw http.ResponseWriter) (jwt.JWT, core.AuthStatus, errors.RichError) {
	logger.Debug("starting getAuthStatus")
	defer logger.Debug("ending getAuthStatus")
	span := trace.SpanFromContext(ctx)
	defer span.End()
	authCookie, err := r.Cookie(constants.LoginCookieName)
	if err != nil {
		// if err == http.ErrNoCookie {
		// according to the r.Cookie method the only possible error is ErrNoCookie,
		// so any error here will be treated as if the cookie was just not there.
		logger.Debug("request did not contain auth cookie")
		return jwt.JWT{}, core.Unauthenticated, nil
	}
	encodedToken := authCookie.Value
	token, rErr := jwt.Decode(encodedToken)
	if err != nil {
		errMsg := "failed to decode jwt from auth cookie jwt"
		logger.Error(errMsg, zap.Reflect("err", rErr))
		apptelemetry.SetSpanOriginalError(&span, rErr, errMsg)
		return jwt.JWT{}, core.Invalid, rErr
	}
	// get / make jwt validator...
	// if valid the we are authenticated
	// handle expired?
	if len(token.Header.KeyID) == 0 {
		err := coreerrors.NewJWTKeyIDMissingError(true)
		errMsg := "jwt key id is required"
		logger.Error(errMsg, zap.Reflect("err", err))
		apptelemetry.SetSpanOriginalError(&span, err, errMsg)
		return jwt.JWT{}, core.Invalid, err
	} else {
		var validator jwt.JWTValidator
		var ok bool
		// here we are making a validator cache key for validating the auth token for the service its self, not for external tokens...
		cachedValidatorKey := "s_" + token.Header.KeyID
		validator, ok = s.validatorCache.GetCachedValidator(cachedValidatorKey)
		if !ok {
			// make the validator and cache it
			jsm, err := s.jsmService.GetJWTSigningMaterialByKeyID(ctx, logger, token.Header.KeyID, "")
			if err != nil {
				errMsg := "jwt signing material for key id not found"
				logger.Error(errMsg, zap.Reflect("err", err))
				apptelemetry.SetSpanError(&span, err, "")
				return jwt.JWT{}, core.Invalid, err
			}
			signer, err := jsm.ToSigner()
			if err != nil {
				errMsg := "failed to create signer from jwt signing material"
				logger.Error(errMsg, zap.Reflect("err", err))
				apptelemetry.SetSpanOriginalError(&span, err, errMsg)
				return jwt.JWT{}, core.Invalid, err
			}
			// validator, err = jwt.NewJWTValidator(jwt.JWTValidatorOptions{})
			validator, err = s.jwtValidatorFactory.NewJWTValidatorWithSigner(token.Header.KeyID, signer)
			if err != nil {
				return jwt.JWT{}, core.Invalid, err
			}
			s.validatorCache.CacheValidator(cachedValidatorKey, validator, cachedJWTValidatorDuration)
		}
		valid, err := validator.ValidateSignature(token.Header.Algorithm, encodedToken)
		if err != nil {
			return jwt.JWT{}, core.Invalid, err
		}
		if !valid {
			return jwt.JWT{}, core.Invalid, nil
		}
		errors, valid := validator.ValidateClaims(token.Claims)
		if !valid {
			logger.Warn("jwt token is not valid")
			for _, e := range errors {
				switch e.GetErrorCode() {
				case coreerrors.ErrCodeExpiredToken: // expiration is a special case and the auth status should reflect it.
					return jwt.JWT{}, core.Expired, nil
				}
				logger.Error("validation of token failed with error(s)", zap.Reflect("err", e))
			}
			return jwt.JWT{}, core.Invalid, coreerrors.NewJWTStandardClaimsInvalidError(encodedToken, true)
		}
		return token, core.Authenticated, nil
	}
}
