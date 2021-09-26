package mongo

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/apptelemetry"
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

func (auditLogRepo) GetName() string {
	return "auditLogRepo"
}

func (auditLogRepo) GetType() string {
	return dataSourceType
}

func (ar auditLogRepo) LogMessage(ctx context.Context, message models.AuditLog) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "LogMessage", ar.GetType())
	defer span.End()
	_, err := ar.mongoClient.Database(ar.dbName).Collection(ar.collection).InsertOne(ctx, message)
	if err != nil {
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, rErr)
	}
	span.AddEvent("audit log added")
	return nil
}
