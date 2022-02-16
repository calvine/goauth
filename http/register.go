package http

import (
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/viewmodels"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

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
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		token, err := s.getNewSCRFToken(ctx, logger)
		if err != nil {
			errorMsg := "failed to create new CSRF token"
			apptelemetry.SetSpanError(&span, err, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		templateRenderError := renderTemplate(rw, registerPageTemplate, viewmodels.RegisterTemplateData{
			CSRFToken: token.Value,
		})
		if templateRenderError != nil {
			errorMsg := "failed to render page template"
			logger.Error(errorMsg,
				zap.Reflect("error", templateRenderError),
			)
			apptelemetry.SetSpanError(&span, templateRenderError, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleRegisterPost() http.HandlerFunc {
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
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
		var errorMsg string
		var templateRenderError errors.RichError
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
		err = s.userService.RegisterUserAndPrimaryContact(ctx, logger, contactType, principal, s.serviceName, "user registration page")
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
			// on success code here...
			accountRegisteredData := viewmodels.AccountRegisteredTemplateData{
				ContactType:      contactType,
				ContactPrincipal: principal,
			}
			templateRenderError = renderTemplate(rw, accountRegisteredPageTemplate, accountRegisteredData)
			if templateRenderError != nil {
				errorMsg = "failed to render template with data provided"
				logger.Error(errorMsg, zap.Reflect("error", templateRenderError), zap.Any("templateData", accountRegisteredData))
				apptelemetry.SetSpanError(&span, templateRenderError, errorMsg)
				// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
				redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			}
			return
		}
	RenderTemplateWithError: // We should only land here if an error occurred that should be forwarded to the client
		token, err := s.getNewSCRFToken(ctx, logger)
		if err != nil {
			errorMsg := "failed to create new CSRF token"
			apptelemetry.SetSpanError(&span, err, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}

		templateData := viewmodels.RegisterTemplateData{
			CSRFToken: token.Value,
			ErrorMsg:  errorMsg,
			Principal: principal,
		}
		templateRenderError = renderTemplate(rw, registerPageTemplate, templateData) //registerPageTemplate.Execute(rw, templateData)
		if templateRenderError != nil {
			errorMsg = "failed to render template with data provided"
			logger.Error(errorMsg, zap.Reflect("error", templateRenderError), zap.Any("templateData", templateData))
			apptelemetry.SetSpanError(&span, templateRenderError, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			redirectToErrorPage(rw, r, errorMsg, http.StatusInternalServerError)
			return
		}
	}
}
