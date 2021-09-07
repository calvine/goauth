package mongo

import (
	"context"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type auditLogRepo struct {
	mongoClient *mongo.Client
	dbName      string
	collection  string
}

func NewAuditLogRepo(client *mongo.Client) *auditLogRepo {
	return &auditLogRepo{client, DB_NAME, AUDITLOG_COLLECTION}
}

func (ar *auditLogRepo) LogMessage(ctx context.Context, message models.AuditLog) errors.RichError {
	_, err := ar.mongoClient.Database(ar.dbName).Collection(ar.collection).InsertOne(ctx, message)
	if err != nil {
		return coreerrors.NewRepoQueryFailedError(err, true)
	}
	return nil
}
