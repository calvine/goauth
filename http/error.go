package http

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/goauth/http/internal/viewmodels"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	// TODO: need to make this a parameter, or a small service in project to return support contact?
	defaultSupportContact = "make_me_configurable@email.com"
)

// var (
// 	errorPageParseSync   sync.Once
// 	errorPageTemplateErr errors.RichError
// )

func (s *server) handleErrorGet() http.HandlerFunc {
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
			if errorPageTemplate == nil {
				errorPageTemplate, templateErr = parseTemplateFromEmbedFS(errorPageTemplatePath, errorPageName, s.templateFS)
			}
		})
		if templateErr != nil {
			errorMsg := "initial parsing of template failed"
			logger.Error(errorMsg, zap.Reflect("error", templateErr))
			apptelemetry.SetSpanError(&span, templateErr, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			http.Error(rw, "RENDER ERROR FAILED!", http.StatusInternalServerError)
			return
		}
		q := r.URL.Query()
		errorMsg := q.Get("errormsg")
		requestID := q.Get("requestid")
		errorCode := q.Get("errorcode")

		templateData := viewmodels.ErrorTemplateData{
			ErrorCode:      errorCode,
			ErrorMessage:   errorMsg,
			RequestID:      requestID,
			SupportContact: defaultSupportContact,
		}

		templateRenderError := renderTemplate(rw, errorPageTemplate, templateData)
		if templateRenderError != nil {
			errorMsg := "failed to render page template"
			logger.Error(errorMsg,
				zap.Reflect("error", templateRenderError),
			)
			apptelemetry.SetSpanError(&span, templateRenderError, errorMsg)
			// s.renderErrorPage(ctx, logger, rw, errorMsg, http.StatusInternalServerError)
			http.Error(rw, "RENDER ERROR FAILED!", http.StatusInternalServerError)
			return
		}
	}
}

func redirectToErrorPage(rw http.ResponseWriter, r *http.Request, errorMsg string, errorCode int) {
	ctx := r.Context()
	requestID := ctxpropagation.GetRequestIDFromContext(ctx)
	queryValues := url.Values{}
	queryValues.Add("errormsg", errorMsg)
	queryValues.Add("errorcode", strconv.Itoa(errorCode))
	queryValues.Add("requestid", requestID)
	url := "/error?" + queryValues.Encode()
	http.Redirect(rw, r, url, http.StatusFound)
}

// func (s *server) renderErrorPage(ctx context.Context, logger *zap.Logger, rw http.ResponseWriter, errorMessage string, httpErrorCode int) {
// 	errorPageParseSync.Do(func() {
// 		if errorPageTemplate == nil {
// 			errorPageTemplate, errorPageTemplateErr = parseTemplateFromEmbedFS(errorPageTemplatePath, errorPageName, s.templateFS)
// 		}
// 	})
// 	if errorPageTemplateErr != nil {
// 		logger.Error("there was an error parsing the error page template", zap.Reflect("error", errorPageTemplateErr))
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		rw.Write([]byte("FAILED TO PARSE ERROR PAGE TEMAPLTE!"))
// 		return
// 	}
// 	if httpErrorCode > 0 {
// 		rw.WriteHeader(httpErrorCode)
// 	}
// 	pageData := viewmodels.ErrorTemplateData{
// 		ErrorMessage:   errorMessage,
// 		RequestID:      ctxpropagation.GetRequestIDFromContext(ctx),
// 		SupportContact: defaultSupportContact,
// 	}
// 	err := errorPageTemplate.Execute(rw, pageData)
// 	if err != nil {
// 		rErr := coreerrors.NewFailedTemplateRenderError("error", err, true)
// 		logger.Error("there was an error rendering the error page template", zap.Reflect("error", rErr))
// 		rw.Write([]byte("UNEXPECTED ERROR OCCURRED..."))
// 	}
// }
