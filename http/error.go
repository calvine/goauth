package http

import (
	"context"
	"net/http"
	"sync"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/utilities/ctxpropagation"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

type errorPageData struct {
	ErrorMessage   string
	RequestID      string
	SupportContact string
}

var (
	errorPageParseSync   sync.Once
	errorPageTemplateErr errors.RichError
)

func (hh *server) renderErrorPage(ctx context.Context, logger *zap.Logger, rw http.ResponseWriter, errorMessage string, httpErrorCode int) {
	errorPageParseSync.Do(func() {
		if errorPageTemplate == nil {
			errorPageTemplate, errorPageTemplateErr = parseTemplateFromEmbedFS(errorPageTemplatePath, errorPageName, hh.templateFS)
		}
	})
	if errorPageTemplate != nil {
		logger.Error("there was an error parsing the error page template", zap.Reflect("error", errorPageTemplate))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("FAILED TO PARSE ERROR PAGE TEMAPLTE!"))
		return
	}
	if httpErrorCode > 0 {
		rw.WriteHeader(httpErrorCode)
	}
	pageData := errorPageData{
		ErrorMessage:   errorMessage,
		RequestID:      ctxpropagation.GetRequestIDFromContext(ctx),
		SupportContact: "make_me_configurable@email.com",
	}
	err := errorPageTemplate.Execute(rw, pageData)
	if err != nil {
		rErr := coreerrors.NewFailedTemplateRenderError("error", err, true)
		logger.Error("there was an error rendering the error page template", zap.Reflect("error", rErr))
		rw.Write([]byte("UNEXPECTED ERROR OCCURRED..."))
	}
}
