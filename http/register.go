package http

import (
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type registerRequestData struct {
	Principal       string
	CSRFToken       string
	ErrorMsg        string
	HasErrorMessage bool
}

func (s *server) handleRegisterGet() http.HandlerFunc {
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
			if registerPageTemplate == nil {
				registerPageTemplate, templateErr = parseTemplateFromEmbedFS(registerPageTemplatePath, registerPageName, s.templateFS)
			}
		})
		if templateErr != nil {
			errorMsg := "initial parsing of template failed"
			logger.Error(errorMsg, zap.Reflect("error", templateErr))
			apptelemetry.SetSpanError(&span, templateErr, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		// TODO: make CSRF token life span configurable
		token, err := models.NewToken("", models.TokenTypeCSRF, constants.Default_CSRF_Token_Duration)
		if err != nil {
			errorMsg := "failed to create new csrf token"
			logger.Error(errorMsg,
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		err = s.tokenService.PutToken(ctx, logger, token)
		if err != nil {
			errorMsg := "failed to store new CSRF token"
			logger.Error(errorMsg,
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
			span.RecordError(err)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		templateRenderError := registerPageTemplate.Execute(rw, registerRequestData{
			CSRFToken:       token.Value,
			HasErrorMessage: false,
		})
		if templateRenderError != nil {
			// TODO: redirect to error page to avoid partial rendered page showing
			errorMsg := "failed to render page template"
			err = coreerrors.NewFailedTemplateRenderError(registerPageName, templateRenderError, true)
			logger.Error(errorMsg,
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleRegisterPost() http.HandlerFunc {
	var (
		once        sync.Once
		templateErr errors.RichError
	)
	type accountRegisteredTemplateData struct {
		ContactType      string
		ContactPrincipal string
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()
		once.Do(func() {
			if registerPageTemplate == nil {
				registerPageTemplate, templateErr = parseTemplateFromEmbedFS(registerPageTemplatePath, registerPageName, s.templateFS)
			}
			if templateErr != nil {
				return
			}
			if accountRegisteredPageTemplate == nil {
				accountRegisteredPageTemplate, templateErr = parseTemplateFromEmbedFS(accountRegisteredPageTemplatePath, accountRegisteredPageName, s.templateFS)
			}
		})
		if templateErr != nil {
			errorMsg := "initial parsing of template failed"
			logger.Error(errorMsg, zap.Reflect("error", templateErr))
			apptelemetry.SetSpanError(&span, templateErr, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		var errorMsg string
		var templateRenderError error
		// TODO: make this a param from the form
		contactType := core.CONTACT_TYPE_EMAIL
		principal := r.FormValue("principal")
		csrfTokenValue := r.FormValue("csrf_token")
		_, err := s.tokenService.GetToken(ctx, logger, csrfTokenValue, models.TokenTypeCSRF)
		if err != nil {
			errorMsg = "failed to retreive csrf token"
			logger.Error(errorMsg,
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
			goto RenderTemplateWithError // I know ewww its a GOTO. I may come back and change how this is being done, possibly extract the render tempalte with error into a function...
		}
		// get principal from request
		err = s.userService.RegisterUserAndPrimaryContact(ctx, logger, contactType, principal, "user registration page")
		if err != nil {
			switch err.GetErrorCode() {
			case coreerrors.ErrCodeInvalidContactPrincipal:
				rw.WriteHeader(http.StatusBadRequest)
				errorMsg = "contact provided is invalid"
			case coreerrors.ErrCodeInvalidContactType:
				rw.WriteHeader(http.StatusBadRequest)
				errorMsg = "contact type provided is invalid"
			case coreerrors.ErrCodeRegistrationContactAlreadyConfirmed:
				// https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.3 seems most appropriate...
				rw.WriteHeader(http.StatusForbidden)
				errorMsg = "contact provided has already been registered"
			default:
				rw.WriteHeader(http.StatusInternalServerError)
				errorMsg = "An error occurred please try again"
			}
			logger.Error("failed to register user with contact provided",
				zap.String("reason", errorMsg),
				zap.String("contactPrincipal", principal),
				zap.String("contactType", contactType),
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, errorMsg)
		} else {
			// FIXME: redirect to page rather than rendre template...
			// on success code here...
			accountRegisteredData := accountRegisteredTemplateData{
				ContactType:      contactType,
				ContactPrincipal: principal,
			}
			templateRenderError = accountRegisteredPageTemplate.Execute(rw, accountRegisteredData)
			if templateRenderError != nil {
				// FIXME: redirect to error page to avoid partial rendered page showing
				errorMsg = "failed to render template with data provided"
				err = coreerrors.NewFailedTemplateRenderError(accountRegisteredPageName, templateRenderError, true)
				logger.Error(errorMsg, zap.Reflect("error", err), zap.Any("templateData", accountRegisteredData))
				apptelemetry.SetSpanError(&span, err, errorMsg)
				s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			}
			return
		}
	RenderTemplateWithError: // We should only land here
		// TODO: make CSRF token life span configurable
		token, err := models.NewToken("", models.TokenTypeCSRF, constants.Default_CSRF_Token_Duration)
		if err != nil {
			errorMsg = "failed to create new CSRF token for page"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
		err = s.tokenService.PutToken(ctx, logger, token)
		if err != nil {
			errorMsg = "failed to store new CSRF token for page"
			logger.Error(errorMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}

		templateData := registerRequestData{
			CSRFToken:       token.Value,
			HasErrorMessage: errorMsg != "",
			ErrorMsg:        errorMsg,
			Principal:       principal,
		}
		templateRenderError = registerPageTemplate.Execute(rw, templateData)
		if templateRenderError != nil {
			// FIXME: redirect to error page to avoid partial rendered page showing
			errorMsg = "failed to render template with data provided"
			err = coreerrors.NewFailedTemplateRenderError(registerPageName, templateRenderError, true)
			logger.Error(errorMsg, zap.Reflect("error", err), zap.Any("templateData", templateData))
			apptelemetry.SetSpanError(&span, err, errorMsg)
			s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			return
		}
	}
}
