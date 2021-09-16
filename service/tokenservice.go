package service

import (
	"context"
	"fmt"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

type tokenService struct {
	tokenRepo repo.TokenRepo
}

func NewTokenService(tokenRepo repo.TokenRepo) services.TokenService {
	return &tokenService{tokenRepo}
}

func (ts tokenService) GetToken(ctx context.Context, tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError) {
	token, err := ts.tokenRepo.GetToken(tokenValue)
	if err != nil {
		return token, err
	}
	now := time.Now().UTC()
	if now.After(token.Expiration) {
		// TODO: do we want to delete the token incase the native store does not support auto delete on TTL like redis?
		return models.Token{}, coreerrors.NewExpiredTokenError(tokenValue, token.TokenType.String(), token.Expiration, true)
	} else if token.TokenType != expectedTokenType {
		// TODO: Audit log this
		return models.Token{}, coreerrors.NewWrongTokenTypeError(token.Value, token.TokenType.String(), expectedTokenType.String(), true)
	}
	return token, nil
}

func (ts tokenService) PutToken(ctx context.Context, token models.Token) errors.RichError {
	tokenErrorsMap := make(map[string]interface{})
	if token.Value == "" {
		// token value must be populated
		tokenErrorsMap["value"] = "token valus is empty"
	} else if token.TokenType == models.TokenTypeInvalid {
		// cannot add invalid token
		tokenErrorsMap["tokenType"] = "token type is invalid"
	} else if token.Expiration.Before(time.Now().UTC()) {
		// cannot save a token that is already expired
		tokenErrorsMap["expiration"] = fmt.Sprintf("token is expired: %s", token.Expiration.String())
	}
	if len(tokenErrorsMap) > 0 {
		return coreerrors.NewMalfomedTokenError(tokenErrorsMap, true)
	}
	err := ts.tokenRepo.PutToken(token)
	return err
}

func (ts tokenService) DeleteToken(ctx context.Context, tokenValue string) errors.RichError {
	err := ts.tokenRepo.DeleteToken(tokenValue)
	return err
}
