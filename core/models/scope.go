package models

type Scope struct {
	ID            string `bson:"-"`
	ApplicationID string `bson:"-"`
	Name          string `bson:"name"`
	Description   string `bson:"description"`
	auditable
}
