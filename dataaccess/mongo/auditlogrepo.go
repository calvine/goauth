package mongo

import (
	"context"

	"github.com/calvine/goauth/core/models"
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

func (ar *auditLogRepo) LogMessage(ctx context.Context, message models.AuditLog) error {
	_, err := ar.mongoClient.Database(ar.dbName).Collection(ar.collection).InsertOne(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
