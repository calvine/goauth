package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	repomodels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var repoJSM repomodels.RepoJWTSigningMaterial
	oid, err := primitive.ObjectIDFromHex(keyID)
	if err != nil {
		rErr := coreerrors.NewFailedToParseObjectIDError(keyID, err, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), keyID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return repoJSM.ToCoreJWTSigningMaterial(), rErr
	}
	filter := bson.M{"_id": oid}
	err = jsm.mongoClient.Database(jsm.dbName).Collection(jsm.collectionName).FindOne(ctx, filter).Decode(&repoJSM)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			rErr := coreerrors.NewNoJWTSigningMaterialFoundError(keyID, true)
			evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), keyID)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return models.JWTSigningMaterial{}, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return models.JWTSigningMaterial{}, rErr
	}
	jwtSigningMaterial := repoJSM.ToCoreJWTSigningMaterial()
	span.AddEvent("jwt signing material retreived")
	return jwtSigningMaterial, nil
}

func (jsm jwtSigningMaterialRepo) AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "AddJWTSigningMaterial", jsm.GetType())
	defer span.End()
	jwtSigningMaterial.AuditData.CreatedByID = createdBy
	jwtSigningMaterial.AuditData.CreatedOnDate = time.Now().UTC()
	result, err := jsm.mongoClient.Database(jsm.dbName).Collection(jsm.collectionName).InsertOne(ctx, jwtSigningMaterial, nil)
	if err != nil {
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	oid := result.InsertedID.(primitive.ObjectID)
	// oid, ok := result.InsertedID.(primitive.ObjectID)
	// if !ok {
	// 	return mongoerrors.NewMongoFailedToParseObjectID(result.InsertedID, true)
	// }
	jwtSigningMaterial.KeyID = oid.Hex()
	span.AddEvent("jwt signing material added")
	return nil
}
