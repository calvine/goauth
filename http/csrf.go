package http

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/http/internal/constants"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// TODO: make CSRF token life span configurable
func (s *server) getNewSCRFToken(ctx context.Context, logger *zap.Logger) (models.Token, errors.RichError) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	token, err := models.NewToken("", models.TokenTypeCSRF, constants.Default_CSRF_Token_Duration)
	if err != nil {
		errorMsg := "failed to create new CSRF token"
		logger.Error(errorMsg, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, errorMsg)
		return models.Token{}, err
	}
	err = s.tokenService.PutToken(ctx, logger, token)
	if err != nil {
		errorMsg := "failed to store new CSRF token"
		logger.Error(errorMsg, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, errorMsg)
		return models.Token{}, err
	}
	return token, err
}

func (s *server) retreiveCSRFToken(ctx context.Context, logger *zap.Logger, value string) (models.Token, errors.RichError) {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	token, err := s.tokenService.GetToken(ctx, logger, value, models.TokenTypeCSRF)
	if err != nil {
		errorMsg := "failed to retreive CSRF token"
		logger.Error(errorMsg, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, errorMsg)
		return models.Token{}, err
	}
	err = s.tokenService.DeleteToken(ctx, logger, value)
	if err != nil {
		warnMsg := "failed to delete CSRF token"
		logger.Warn(warnMsg, zap.Reflect("error", err))
		// apptelemetry.SetSpanError(&span, err, errorMsg)
		// I dont think we want to return an error here...
	}
	return token, nil
}
