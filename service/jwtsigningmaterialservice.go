package service

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
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

func (jsms jwtSigningMaterialService) AddJWTSigningMaterial(ctx context.Context, logger *zap.Logger, jsm *models.JWTSigningMaterial, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "AddJWTSigningMaterial")
	defer span.End()
	if jsm.KeyID == "" {
		err := coreerrors.NewJWTSigningMaterialKeyIDMissingError(true)
		evtString := "jwt signing material key id is empty"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	if jsm.AlgorithmType == "" {
		err := coreerrors.NewJWTSigningMaterialAlgorithmTypeMissingError(true)
		evtString := "jwt signing material algorithm type is empty"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// check for existing record with the key id provided
	exists := true
	_, err := jsms.jsmRepo.GetJWTSigningMaterialByKeyID(ctx, jsm.KeyID)
	if err != nil {
		if err.GetErrorCode() == coreerrors.ErrCodeNoJWTSigningMaterialFound {
			evtString := "no existing jwt signing material with given key id found"
			span.AddEvent(evtString)
			logger.Info(evtString)
			exists = false
		} else {
			logger.Error("jsmRepo.AddJWTSigningMaterial call failed", zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, "")
			return err
		}
	}
	if exists {
		err = coreerrors.NewJWTSigningMaterialKeyIDNotUniqueError(jsm.KeyID, true)
		evtString := "jwt signing material key id already exists"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	err = jsms.jsmRepo.AddJWTSigningMaterial(ctx, jsm, initiator)
	return err
}

func (jsms jwtSigningMaterialService) GetJWTSigningMaterialByKeyID(ctx context.Context, logger *zap.Logger, keyID string, initiator string) (models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "GetJWTSigningMaterialByKeyID")
	defer span.End()

	jsm, err := jsms.jsmRepo.GetJWTSigningMaterialByKeyID(ctx, keyID)
	if err != nil {
		logger.Error("jsmRepo.GetJWTSigningMaterialByKeyID call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.JWTSigningMaterial{}, err
	}
	if jsm.IsExpired() {
		err := coreerrors.NewJWTSigningMaterialExpiredError(keyID, true)
		evtString := "jwt signing material expired"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.JWTSigningMaterial{}, err
	}

	return jsm, nil
}

func (jsms jwtSigningMaterialService) GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, logger *zap.Logger, algorithmType jwt.JWTSigningAlgorithmFamily, initiator string) ([]models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, jsms.GetName(), "GetValidJWTSigningMaterialByAlgorithmType")
	defer span.End()

	results, err := jsms.jsmRepo.GetValidJWTSigningMaterialByAlgorithmType(ctx, algorithmType)
	if err != nil {
		logger.Error("jsmRepo.GetValidJWTSigningMaterialByAlgorithmType call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return nil, err
	}

	return results, nil
}
