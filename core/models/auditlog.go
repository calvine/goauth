package models

import "time"

const (
	AssetType_User        = "user"
	AssetType_Application = "application"
)

type AuditLog struct {
	ID           string                 `bson:"id"`
	Message      string                 `bson:"message"`
	Code         string                 `bson:"code"`
	AssetType    string                 `bson:"assetType"`
	AssetID      string                 `bson:"assetId"`
	AuditLogDate time.Time              `bson:"auditLogDate"`
	Data         map[string]interface{} `bson:"data"`
}
