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

type localTokenRepo struct {
	tokenMap map[string]models.Token
}

func NewMemoryTokenRepo() repo.TokenRepo {
	tokenMap := make(map[string]models.Token)
	return &localTokenRepo{tokenMap}
}

func (ltr *localTokenRepo) GetToken(ctx context.Context, tokenValue string) (models.Token, errors.RichError) {
	token, ok := ltr.tokenMap[tokenValue]
	if !ok {
		return token, coreerrors.NewTokenNotFoundError(tokenValue, true)
	}
	return token, nil
}

func (ltr *localTokenRepo) PutToken(ctx context.Context, token models.Token) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer("tokenRepo").Start(ctx, "PutToken")
	span.SetAttributes(attribute.String("db", "memory"))
	defer span.End()
	ltr.tokenMap[token.Value] = token
	return nil
}

func (ltr *localTokenRepo) DeleteToken(ctx context.Context, tokenValue string) errors.RichError {
	delete(ltr.tokenMap, tokenValue)
	return nil
}
