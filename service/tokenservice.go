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
	"go.uber.org/zap"
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

func (ts tokenService) GetToken(ctx context.Context, logger *zap.Logger, tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "GetToken")
	defer span.End()
	token, err := ts.tokenRepo.GetToken(ctx, tokenValue)
	if err != nil {
		apptelemetry.SetSpanError(&span, err, "")
		logger.Error("tokenRepo.GetToken call failed", zap.Reflect("error", err))
		return token, err
	}
	span.AddEvent("token retreived from tokenRepo")
	if token.IsExpired() {
		// TODO: do we want to delete the token incase the native store does not support auto delete on TTL like redis?
		evtString := fmt.Sprintf("token expired on %s", token.Expiration.UTC().String())
		err := coreerrors.NewExpiredTokenError(tokenValue, token.TokenType.String(), token.Expiration, true)
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.Token{}, err
	} else if token.TokenType != expectedTokenType {
		// TODO: Audit log this
		evtString := fmt.Sprintf("token type %s does not match expected type %s", token.TokenType.String(), expectedTokenType.String())
		err := coreerrors.NewWrongTokenTypeError(token.Value, token.TokenType.String(), expectedTokenType.String(), true)
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.Token{}, err
	}
	span.AddEvent("token retrevied")
	return token, nil
}

func (ts tokenService) PutToken(ctx context.Context, logger *zap.Logger, token models.Token) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "PutToken")
	defer span.End()
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
		err := coreerrors.NewMalfomedTokenError(tokenErrorsMap, true)
		evtString := "token validation failed"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("token validated")
	err := ts.tokenRepo.PutToken(ctx, token)
	if err != nil {
		logger.Error("tokenRepo.PutToken call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("token stored")
	return nil
}

func (ts tokenService) DeleteToken(ctx context.Context, logger *zap.Logger, tokenValue string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ts.GetName(), "DeleteToken")
	defer span.End()
	err := ts.tokenRepo.DeleteToken(ctx, tokenValue)
	if err != nil {
		logger.Error("tokenRepo.DeleteToken call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	evtString := fmt.Sprintf("token deleted: %s", tokenValue)
	span.AddEvent(evtString)
	return nil
}
