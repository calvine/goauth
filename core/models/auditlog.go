package models

import "time"

const (
	AssetType_User        = "user"
	AssetType_Application = "application"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFORMATION
	WARNING
	ERROR
	CRITICAL
)

func (ll LogLevel) ToString() string {
	switch ll {
	case DEBUG:
		return "DBUG"
	case INFORMATION:
		return "INFO"
	case WARNING:
		return "WARN"
	case ERROR:
		return "EROR"
	case CRITICAL:
		return "CRIT"
	default:
		// should not be possible...
		return "INVALID"
	}
}

type AuditLog struct {
	ID           string                 `bson:"id"`
	Message      string                 `bson:"message"`
	Code         string                 `bson:"code"`
	AssetType    string                 `bson:"assetType"`
	AssetID      string                 `bson:"assetId"`
	AuditLogDate time.Time              `bson:"auditLogDate"`
	Data         map[string]interface{} `bson:"data"`
}
