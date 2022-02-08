package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

type jwtSigningMaterialRepo struct {
	material *map[string]models.JWTSigningMaterial
}

func NewMemoryJWTSigningMaterialRepo() repo.JWTSigningMaterialRepo {
	material := make(map[string]models.JWTSigningMaterial)
	return jwtSigningMaterialRepo{
		material: &material,
	}
}

func (jwtSigningMaterialRepo) GetName() string {
	return "jwtSigningMaterialRepo"
}

func (jwtSigningMaterialRepo) GetType() string {
	return dataSourceType
}

func (jsm jwtSigningMaterialRepo) GetJWTSigningMaterialByKeyID(ctx context.Context, keyID string) (models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "GetJWTSigningMaterialByKeyID", jsm.GetType())
	defer span.End()
	material, ok := (*jsm.material)[keyID]
	if !ok {
		fields := map[string]interface{}{
			"keyID": keyID,
		}
		rErr := coreerrors.NewNoJWTSigningMaterialFoundError(fields, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), keyID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return models.JWTSigningMaterial{}, rErr
	}
	return material, nil
}

func (jsm jwtSigningMaterialRepo) GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, algorithmType jwt.JWTSigningAlgorithmFamily) ([]models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "GetJWTSigningMaterialByKeyID", jsm.GetType())
	defer span.End()
	results := make([]models.JWTSigningMaterial, 0, 5)
	for _, sm := range *jsm.material {
		if !sm.Disabled && !sm.IsExpired() && sm.AlgorithmType == algorithmType {
			results = append(results, sm)
		}
	}
	// if len(results) == 0 {
	// 	fields := map[string]interface{}{
	// 		"algorithmType": algorithmType,
	// 	}
	// 	rErr := coreerrors.NewNoJWTSigningMaterialFoundError(fields, true)
	// 	evtString := fmt.Sprintf("no jwt signing material found for algorithm type: %s", algorithmType)
	// 	apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
	// 	return nil, rErr
	// }
	return results, nil
}

func (jsm jwtSigningMaterialRepo) AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "AddJWTSigningMaterial", jsm.GetType())
	defer span.End()
	if jwtSigningMaterial.ID == "" {
		// if id is not set lets set one
		jwtSigningMaterial.ID = uuid.New().String()
	}
	/*
		This logic should be in the service level, not in the repo...
	*/
	// _, alreadyExists := (*jsm.material)[jwtSigningMaterial.KeyID]
	// if alreadyExists {
	// 	rErr := coreerrors.NewJWTSigningMaterialKeyIDNotUniqueError(jwtSigningMaterial.KeyID, true)
	// 	evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), jwtSigningMaterial.KeyID)
	// 	apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
	// 	return rErr
	// }
	jwtSigningMaterial.AuditData.CreatedByID = createdBy
	jwtSigningMaterial.AuditData.CreatedOnDate = time.Now().UTC()
	(*jsm.material)[jwtSigningMaterial.KeyID] = *jwtSigningMaterial
	return nil

}
