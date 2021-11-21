package http

import (
	"html/template"
	"net/http"
	"sync"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/constants"
	"go.opentelemetry.io/otel/trace"
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
		// TODO: make CSRF token life span configurable
		ctx := r.Context()
		logger := ctxpropagation.GetLoggerFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		// TODO: make this a param from the form
		contactType := core.CONTACT_TYPE_EMAIL
		// get principal from request
		principal := r.FormValue("principal")
		err := s.userService.RegisterUserAndPrimaryContact(ctx, logger, contactType, principal, "user registration page")
		if err != nil {

		}
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
