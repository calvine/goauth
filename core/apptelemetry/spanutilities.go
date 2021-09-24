package apptelemetry

import (
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func SetSpanError(span *trace.Span, err errors.RichError) {
	(*span).SetAttributes(
		attribute.String("errorCode", err.GetErrorCode()),
		// attribute.String("errorMessage", err.GetErrorMessage()),
	)
	(*span).SetStatus(codes.Error, err.ToString(errors.ShortDetailedOutput))
}
