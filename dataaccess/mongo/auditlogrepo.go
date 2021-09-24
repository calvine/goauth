package mongo

import (
	"context"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type auditLogRepo struct {
	mongoClient *mongo.Client
	dbName      string
	collection  string
}

func NewAuditLogRepo(client *mongo.Client) *auditLogRepo {
	return &auditLogRepo{client, DB_NAME, AUDITLOG_COLLECTION}
}

func (auditLogRepo) GetName() string {
	return "auditLogRepo"
}

func (auditLogRepo) GetType() string {
	return dataSourceType
}

func (ar auditLogRepo) LogMessage(ctx context.Context, message models.AuditLog) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ar.GetName()).Start(ctx, "LogMessage")
	span.SetAttributes(attribute.String("db", ar.GetType()))
	defer span.End()
	_, err := ar.mongoClient.Database(ar.dbName).Collection(ar.collection).InsertOne(ctx, message)
	if err != nil {
		return coreerrors.NewRepoQueryFailedError(err, true)
	}
	return nil
}
