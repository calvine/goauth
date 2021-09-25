package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
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

func (tokenService) GetName() string {
	return "tokenService"
}

func (ts tokenService) GetToken(ctx context.Context, tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "GetToken")
	defer span.End()
	token, err := ts.tokenRepo.GetToken(ctx, tokenValue)
	if err != nil {
		evtString := fmt.Sprintf("failed to retreive token from token repo: %s - %s", tokenValue, expectedTokenType.String())
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return token, err
	}
	span.AddEvent("token retreived from tokenRepo")
	now := time.Now().UTC()
	if now.After(token.Expiration) {
		// TODO: do we want to delete the token incase the native store does not support auto delete on TTL like redis?
		evtString := fmt.Sprintf("token expired on %s", token.Expiration.UTC().String())
		span.AddEvent(evtString)
		err := coreerrors.NewExpiredTokenError(tokenValue, token.TokenType.String(), token.Expiration, true)
		apptelemetry.SetSpanError(&span, err)
		return models.Token{}, err
	} else if token.TokenType != expectedTokenType {
		// TODO: Audit log this
		evtString := fmt.Sprintf("token type %s does not match expected type %s", token.TokenType.String(), expectedTokenType.String())
		span.AddEvent(evtString)
		err := coreerrors.NewWrongTokenTypeError(token.Value, token.TokenType.String(), expectedTokenType.String(), true)
		apptelemetry.SetSpanError(&span, err)
		return models.Token{}, err
	}
	return token, nil
}

func (ts tokenService) PutToken(ctx context.Context, token models.Token) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "PutToken")
	defer span.End()
	tokenErrorsMap := make(map[string]interface{})
	if token.Value == "" {
		// token value must be populated
		span.AddEvent("token value is empty string")
		tokenErrorsMap["value"] = "token valus is empty"
	} else if token.TokenType == models.TokenTypeInvalid {
		// cannot add invalid token
		evtString := fmt.Sprintf("token type is invalid: %s", token.TokenType.String())
		span.AddEvent(evtString)
		tokenErrorsMap["tokenType"] = "token type is invalid"
	} else if token.Expiration.Before(time.Now().UTC()) {
		// cannot save a token that is already expired
		evtString := fmt.Sprintf("token expiration is in the past: %s", token.Expiration.UTC().String())
		span.AddEvent(evtString)
		tokenErrorsMap["expiration"] = fmt.Sprintf("token is expired: %s", token.Expiration.String())
	}
	if len(tokenErrorsMap) > 0 {
		span.AddEvent("token validation failed")
		err := coreerrors.NewMalfomedTokenError(tokenErrorsMap, true)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	err := ts.tokenRepo.PutToken(ctx, token)
	if err != nil {
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	return nil
}

func (ts tokenService) DeleteToken(ctx context.Context, tokenValue string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "DeleteToken")
	defer span.End()
	err := ts.tokenRepo.DeleteToken(ctx, tokenValue)
	if err != nil {
		evtString := fmt.Sprintf("failed to delete token: %s", tokenValue)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	evtString := fmt.Sprintf("token deleted: %s", tokenValue)
	span.AddEvent(evtString)
	return nil
}
