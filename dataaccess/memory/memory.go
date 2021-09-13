package memory

import "github.com/calvine/goauth/core/models"

var (
	users     map[string]models.User
	contacts  map[string]models.Contact
	apps      map[string]models.App
	appScopes map[string][]models.Scope
)
