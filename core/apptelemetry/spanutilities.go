package apptelemetry

import (
	"context"

	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	DataSourceTypeKey = "db"
)

func CreateRepoFunctionSpan(ctx context.Context, componentName string, funcName string, dataSourceType string, attributes ...attribute.KeyValue) trace.Span {
	span := CreateFunctionSpan(ctx, componentName, funcName, attributes...)
	span.SetAttributes(attribute.String(DataSourceTypeKey, dataSourceType))
	return span
}

func CreateFunctionSpan(ctx context.Context, componentName string, funcName string, attributes ...attribute.KeyValue) trace.Span {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(componentName).Start(ctx, funcName)
	if len(attributes) > 0 {
		span.SetAttributes(attributes...)
	}
	return span
}
func SetSpanOriginalError(span *trace.Span, err errors.RichError, eventString string) {
	(*span).AddEvent(eventString)
	(*span).SetAttributes(
		attribute.String("errorData", err.Error()),
		// attribute.String("errorMessage", err.GetErrorMessage()),
	)
	(*span).SetStatus(codes.Error, err.GetErrorCode())
}

func SetSpanError(span *trace.Span, err errors.RichError, optionalEventString string) {
	if optionalEventString != "" {
		(*span).AddEvent(optionalEventString)
	}
	(*span).SetStatus(codes.Error, err.GetErrorCode())
}
