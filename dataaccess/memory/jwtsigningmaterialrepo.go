package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
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
		rErr := coreerrors.NewNoJWTSigningMaterialFoundError(keyID, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), keyID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return models.JWTSigningMaterial{}, rErr
	}
	return material, nil
}

func (jsm jwtSigningMaterialRepo) AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "AddJWTSigningMaterial", jsm.GetType())
	defer span.End()
	_, alreadyExists := (*jsm.material)[jwtSigningMaterial.KeyID]
	if alreadyExists {
		rErr := coreerrors.NewJWTSigningMaterialKeyIDNotUniqueError(jwtSigningMaterial.KeyID, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), jwtSigningMaterial.KeyID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	jwtSigningMaterial.AuditData.CreatedByID = createdBy
	jwtSigningMaterial.AuditData.CreatedOnDate = time.Now().UTC()
	(*jsm.material)[jwtSigningMaterial.KeyID] = *jwtSigningMaterial
	return nil

}
