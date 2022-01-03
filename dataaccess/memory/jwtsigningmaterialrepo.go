package memory

import (
	"context"
	"time"

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
	return "userRepo"
}

func (jwtSigningMaterialRepo) GetType() string {
	return dataSourceType
}

func (jsm jwtSigningMaterialRepo) GetJWTSigningMaterialByKeyID(ctx context.Context, keyID string) (models.JWTSigningMaterial, errors.RichError) {
	material, ok := (*jsm.material)[keyID]
	if !ok {
		fields := make(map[string]interface{})
		fields["keyID"] = keyID
		return models.JWTSigningMaterial{}, coreerrors.NewNoJWTSigningMaterialFoundError(keyID, true)
	}
	return material, nil
}

func (jsm jwtSigningMaterialRepo) AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError {
	_, alreadyExists := (*jsm.material)[jwtSigningMaterial.KeyID]
	if alreadyExists {
		return coreerrors.NewJWTSigningMaterialKeyIDNotUniqueError(jwtSigningMaterial.KeyID, true)
	}
	jwtSigningMaterial.AuditData.CreatedByID = createdBy
	jwtSigningMaterial.AuditData.CreatedOnDate = time.Now().UTC()
	(*jsm.material)[jwtSigningMaterial.KeyID] = *jwtSigningMaterial
	return nil

}
