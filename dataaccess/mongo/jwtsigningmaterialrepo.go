package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
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
	filter := bson.M{"keyId": keyID}
	err := jsm.mongoClient.Database(jsm.dbName).Collection(jsm.collectionName).FindOne(ctx, filter).Decode(&repoJSM)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"keyID": keyID,
			}
			rErr := coreerrors.NewNoJWTSigningMaterialFoundError(fields, true)
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

func (jsm jwtSigningMaterialRepo) GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, algorithmType jwt.JWTSingingAlgorithmFamily) ([]models.JWTSigningMaterial, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, jsm.GetName(), "GetValidJWTSigningMaterialByAlgorithmType", jsm.GetType())
	defer span.End()
	repoResult := make([]repomodels.RepoJWTSigningMaterial, 0, 5)
	filter := bson.M{
		"$and": bson.A{
			bson.M{"disabled": false},
			bson.M{
				"$or": bson.A{
					bson.D{
						{
							Key:   "expiration",
							Value: nil,
						},
					},
					bson.D{
						{
							Key: "expiration",
							Value: bson.M{
								"$gt": time.Now().UTC(),
							},
						},
					},
				},
			},
			bson.M{"algorithmType": algorithmType},
		},
	}
	queryResult, err := jsm.mongoClient.Database(jsm.dbName).Collection(jsm.collectionName).Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"algorithmType": algorithmType,
			}
			rErr := coreerrors.NewNoJWTSigningMaterialFoundError(fields, true)
			evtString := fmt.Sprintf("no jwt signing material found for algorithm type: %s", algorithmType)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return nil, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	defer queryResult.Close(ctx)
	err = queryResult.All(ctx, &repoResult)
	if err != nil {
		rErr := coreerrors.NewFailedToDecodeRepoDataError("RepoJWTSigningMaterial", err, true)
		evtString := fmt.Sprintf("failed to decode results from repo query: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	numResults := len(repoResult)
	result := make([]models.JWTSigningMaterial, numResults)
	for i := range repoResult {
		result[i] = models.JWTSigningMaterial(repoResult[i].ToCoreJWTSigningMaterial())
	}
	return result, nil

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
	jwtSigningMaterial.ID = oid.Hex()
	span.AddEvent("jwt signing material added")
	return nil
}
