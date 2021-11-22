package http

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
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
		once             sync.Once
		registerTemplate *template.Template
		templateErr      error
		templatePath     string = "http/templates/register.html.tmpl"
	)
	return func(rw http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			templateFileData, err := s.templateFS.ReadFile(templatePath)
			templateErr = err
			if templateErr == nil {
				registerTemplate, templateErr = template.New("registerPage").Parse(string(templateFileData))
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
		defer span.End()
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
		templateRenderError := registerTemplate.Execute(rw, registerRequestData{
			CSRFToken:       token.Value,
			HasErrorMessage: false,
		})
		if templateRenderError != nil {
			span.RecordError(err)
			err = coreerrors.NewTemplateRenderErrorError(templatePath, templateRenderError, true)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleRegisterPost() http.HandlerFunc {
	var (
		once             sync.Once
		registerTemplate *template.Template
		templateErr      error
		templatePath     string = "http/templates/register.html.tmpl"
	)
	return func(rw http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			templateFileData, err := s.templateFS.ReadFile(templatePath)
			templateErr = err
			if templateErr == nil {
				registerTemplate, templateErr = template.New("registerPage").Parse(string(templateFileData))
			}
		})
		if templateErr != nil {
			http.Error(rw, templateErr.Error(), http.StatusInternalServerError)
			return
		}
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		defer span.End()
		csrfTokenValue := r.FormValue("principal")
		_, err := s.tokenService.GetToken(ctx, logger, csrfTokenValue, models.TokenTypeCSRF)
		if err != nil {
			logger.Error("failed to retreive csrf token",
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, "")
		}
		// TODO: make this a param from the form
		contactType := core.CONTACT_TYPE_EMAIL
		// get principal from request
		principal := r.FormValue("principal")
		err = s.userService.RegisterUserAndPrimaryContact(ctx, logger, contactType, principal, "user registration page")
		if err != nil {
			logger.Error("failed to register user with contact provided",
				zap.String("contactPrincipal", principal),
				zap.String("contactType", contactType),
				zap.Reflect("error", err),
			)
			apptelemetry.SetSpanError(&span, err, "")
			var errorMsg string
			switch err.GetErrorCode() {
			case coreerrors.ErrCodeInvalidContactPrincipal:
				rw.WriteHeader(http.StatusBadRequest)
				errorMsg = "contact provided is invalid"
				return
			case coreerrors.ErrCodeInvalidContactType:
				rw.WriteHeader(http.StatusBadRequest)
				errorMsg = "contact type provided is invalid"
				return
			case coreerrors.ErrCodeRegistrationContactAlreadyConfirmed:
				// https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.3 seems most appropriate...
				rw.WriteHeader(http.StatusForbidden)
				errorMsg = "contact provided has already been registered"
				return
			default:
				rw.WriteHeader(http.StatusInternalServerError)
				errorMsg = "An error occurred please try again"
			}
			// TODO: make CSRF token life span configurable
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
			templateRenderError := registerTemplate.Execute(rw, registerRequestData{
				CSRFToken:       token.Value,
				HasErrorMessage: errorMsg != "",
				ErrorMsg:        errorMsg,
				Principal:       principal,
			})
			if templateRenderError != nil {
				span.RecordError(err)
				err = coreerrors.NewTemplateRenderErrorError(templatePath, templateRenderError, true)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// on success code here...
		// TOOD: make a registered static page indicating that a notification was sent and that is how to finish registration...
	}
}
