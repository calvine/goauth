package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

type Scope struct {
	ID          string    `bson:"-"`
	AppID       string    `bson:"-"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	AuditData   auditable `bson:",inline"`
}

func NewScope(appID, name, description string) Scope {
	return Scope{
		AppID:       appID,
		Name:        name,
		Description: description,
	}
}

func ValidateScope(includeID bool, scope Scope) errors.RichError {
	fields := make(map[string]interface{})
	if includeID && scope.ID == "" {
		fields["ID"] = "app ID cannot be empty"
	}
	if scope.AppID == "" {
		fields["AppID"] = "app AppID cannot be empty"
	}
	if scope.Name == "" {
		fields["Name"] = "app Name cannot be empty"
	}
	if scope.Description == "" {
		fields["Description"] = "app Description cannot be empty"
	}

	if len(fields) > 0 {
		return coreerrors.NewInvalidScopeCreationError(fields, false)
	}
	return nil
}
