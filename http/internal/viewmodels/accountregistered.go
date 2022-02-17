package viewmodels

import "github.com/calvine/goauth/core"

type AccountRegisteredTemplateData struct {
	ContactType      core.ContactType
	ContactPrincipal string
}
