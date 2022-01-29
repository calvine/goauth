package service

import (
	"context"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

type cachedJSM struct {
	aquired  time.Time
	material models.JWTSigningMaterial
}

type cachedJWTSigningMaterialService struct {
	cacheDuration time.Duration
	keyIDCache    map[string]cachedJSM
	innerService  services.JWTSigningMaterialService
}

func NewCachedJWTSigningMaterialService(s services.JWTSigningMaterialService, cacheDuration time.Duration) services.JWTSigningMaterialService {
	return cachedJWTSigningMaterialService{
		cacheDuration: cacheDuration,
		innerService:  s,
	}
}

func (cachedJWTSigningMaterialService) GetName() string {
	return "cachedJWTSigningMaterialService"
}

func (cjsms cachedJWTSigningMaterialService) AddJWTSigningMaterial(ctx context.Context, logger *zap.Logger, jsm *models.JWTSigningMaterial, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, cjsms.GetName(), "AddJWTSigningMaterial")
	defer span.End()

	return cjsms.innerService.AddJWTSigningMaterial(ctx, logger, jsm, initiator)
}

func (cjsms cachedJWTSigningMaterialService) GetJWTSigningMaterialByKeyID(ctx context.Context, logger *zap.Logger, keyID string, initiator string) (models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, cjsms.GetName(), "GetJWTSigningMaterialByKeyID")
	defer span.End()
	cachedData, ok := cjsms.keyIDCache[keyID]
	if !ok {
		// cache miss
		span.AddEvent("cache miss")
		logger.Info("cache miss for key id", zap.String("keyID", keyID))
		return cjsms.innerService.GetJWTSigningMaterialByKeyID(ctx, logger, keyID, initiator)
	}
	age := time.Since(cachedData.aquired)
	if age > cjsms.cacheDuration {
		// cache expired
		delete(cjsms.keyIDCache, keyID)
		span.AddEvent("cache expired")
		logger.Info("cache expired for key id", zap.String("keyID", keyID))
		return cjsms.innerService.GetJWTSigningMaterialByKeyID(ctx, logger, keyID, initiator)
	}
	// cache hit
	span.AddEvent("cache hit")
	logger.Info("cache hit for key id!", zap.String("keyID", keyID))
	return cachedData.material, nil
}

func (cjsms cachedJWTSigningMaterialService) GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, logger *zap.Logger, algorithmType jwt.JWTSingingAlgorithmFamily, initiator string) ([]models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, cjsms.GetName(), "GetValidJWTSigningMaterialByAlgorithmType")
	defer span.End()
	// no cache here yet...
	return cjsms.innerService.GetValidJWTSigningMaterialByAlgorithmType(ctx, logger, algorithmType, initiator)
}
