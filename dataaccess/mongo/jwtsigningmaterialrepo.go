package mongo

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type jwtSigningMaterialRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewJWTSigningMaterialRepo(client *mongo.Client) repo.JWTSigningMaterialRepo {
	return jwtSigningMaterialRepo{client, DB_NAME, JWT_SIGNING_MATERIAL_COLLECTION}
}

func NewJWTSigningMaterialRepoWithNames(client *mongo.Client, dbName, collectionName string) repo.JWTSigningMaterialRepo {
	return jwtSigningMaterialRepo{client, dbName, collectionName}
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
	// material, ok := (*jsm.material)[keyID]
	// if !ok {
	// 	fields := make(map[string]interface{})
	// 	fields["keyID"] = keyID
	// 	return models.JWTSigningMaterial{}, coreerrors.NewNoJWTSigningMaterialFoundError(keyID, true)
	// }
	// return material, nil
}

func (jsm jwtSigningMaterialRepo) AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "AddJWTSigningMaterial", jsm.GetType())
	defer span.End()
	// _, alreadyExists := (*jsm.material)[jwtSigningMaterial.KeyID]
	// if alreadyExists {
	// 	return coreerrors.NewJWTSigningMaterialKeyIDNotUniqueError(jwtSigningMaterial.KeyID, true)
	// }
	// jwtSigningMaterial.AuditData.CreatedByID = createdBy
	// jwtSigningMaterial.AuditData.CreatedOnDate = time.Now().UTC()
	// (*jsm.material)[jwtSigningMaterial.KeyID] = *jwtSigningMaterial
	// return nil
}
