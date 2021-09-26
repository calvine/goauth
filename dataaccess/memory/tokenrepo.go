package memory

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
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
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ltr.GetName(), "GetToken", ltr.GetType())
	defer span.End()
	token, ok := ltr.tokenMap[tokenValue]
	if !ok {
		evtString := fmt.Sprintf("token not found: %s", tokenValue)
		err := coreerrors.NewTokenNotFoundError(tokenValue, true)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return token, err
	}
	span.AddEvent("token retreived")
	return token, nil
}

func (ltr *tokenRepo) PutToken(ctx context.Context, token models.Token) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ltr.GetName(), "PutToken", ltr.GetType())
	defer span.End()
	ltr.tokenMap[token.Value] = token
	span.AddEvent("token stored")
	return nil
}

func (ltr *tokenRepo) DeleteToken(ctx context.Context, tokenValue string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ltr.GetName(), "DeleteToken", ltr.GetType())
	defer span.End()
	_, ok := ltr.tokenMap[tokenValue]
	if !ok {
		evtString := fmt.Sprintf("token not found: %s", tokenValue)
		err := coreerrors.NewTokenNotFoundError(tokenValue, true)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	delete(ltr.tokenMap, tokenValue)
	span.AddEvent("token deleted")
	return nil
}
