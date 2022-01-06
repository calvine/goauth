package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

type Scope struct {
	ID          string    `bson:"-"`
	AppID       string    `bson:"-"`
	Name        string    `bson:"name"`
	DisplayName string    `bson:"displayName"`
	Description string    `bson:"description"`
	AuditData   auditable `bson:",inline"`
}

func NewScope(appID, name, displayName, description string) Scope {
	return Scope{
		AppID:       appID,
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}
}

func ValidateScope(includeID bool, scope Scope) errors.RichError {
	fields := make(map[string]interface{})
	if includeID && scope.ID == "" {
		fields["ID"] = "ID cannot be empty"
	}
	if scope.AppID == "" {
		fields["AppID"] = "app AppID cannot be empty"
	}
	if scope.Name == "" {
		fields["Name"] = "Name cannot be empty"
	}
	if scope.DisplayName == "" {
		fields["DisplayName"] = "DisplayName cannot be empty"
	}
	if scope.Description == "" {
		fields["Description"] = "Description cannot be empty"
	}

	if len(fields) > 0 {
		return coreerrors.NewInvalidScopeCreationError(fields, false)
	}
	return nil
}
