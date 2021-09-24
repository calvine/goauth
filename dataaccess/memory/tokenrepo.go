package memory

import (
	"context"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tokenRepo struct {
	tokenMap map[string]models.Token
}

func NewMemoryTokenRepo() repo.TokenRepo {
	tokenMap := make(map[string]models.Token)
	return &tokenRepo{tokenMap}
}

func (tokenRepo) GetName() string {
	return "tokenRepo"
}

func (tokenRepo) GetType() string {
	return dataSourceType
}

func (ltr *tokenRepo) GetToken(ctx context.Context, tokenValue string) (models.Token, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ltr.GetName()).Start(ctx, "GetToken")
	span.SetAttributes(attribute.String("db", ltr.GetType()))
	defer span.End()
	token, ok := ltr.tokenMap[tokenValue]
	if !ok {
		return token, coreerrors.NewTokenNotFoundError(tokenValue, true)
	}
	return token, nil
}

func (ltr *tokenRepo) PutToken(ctx context.Context, token models.Token) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ltr.GetName()).Start(ctx, "PutToken")
	span.SetAttributes(attribute.String("db", ltr.GetType()))
	defer span.End()
	ltr.tokenMap[token.Value] = token
	return nil
}

func (ltr *tokenRepo) DeleteToken(ctx context.Context, tokenValue string) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ltr.GetName()).Start(ctx, "DeleteToken")
	span.SetAttributes(attribute.String("db", ltr.GetType()))
	defer span.End()
	delete(ltr.tokenMap, tokenValue)
	return nil
}
