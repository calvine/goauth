package models

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
