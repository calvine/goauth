package service

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	corerepo "github.com/calvine/goauth/core/repositories"
	coreservices "github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

type jwtSigningMaterialService struct {
	jsmRepo corerepo.JWTSigningMaterialRepo
}

func NewJWTSigningMaterialService(jsmRepo corerepo.JWTSigningMaterialRepo) coreservices.JWTSigningMaterialService {
	return jwtSigningMaterialService{
		jsmRepo: jsmRepo,
	}
}

func (jwtSigningMaterialService) GetName() string {
	return "jwtSigningMaterialService"
}

func (jsms jwtSigningMaterialService) AddJWTSigningMaterial(ctx context.Context, logger *zap.Logger, jsm *models.JWTSigningMaterial) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "AddJWTSigningMaterial")
	defer span.End()

	return coreerrors.NewNotImplementedError(true)
}

func (jsms jwtSigningMaterialService) GetJWTSigningMaterialByKeyID(ctx context.Context, logger *zap.Logger, keyID string) (models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "GetJWTSigningMaterialByKeyID")
	defer span.End()

	return models.JWTSigningMaterial{}, coreerrors.NewNotImplementedError(true)
}

func (jsms jwtSigningMaterialService) GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, logger *zap.Logger, algorithmType string) ([]models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "GetValidJWTSigningMaterialByAlgorithmType")
	defer span.End()

	return nil, coreerrors.NewNotImplementedError(true)
}
