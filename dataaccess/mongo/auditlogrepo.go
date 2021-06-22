package mongo

import (
	"github.com/calvine/goauth/models/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type auditLogRepo struct {
	mongoClient *mongo.Client
	dbName      string
	collection  string
}

func NewAuditLogRepo(client *mongo.Client) *auditLogRepo {
	return &auditLogRepo{client, "", ""}
}

type AuditLogRepo interface {
	LogMessage(message core.AuditLog) error
}
