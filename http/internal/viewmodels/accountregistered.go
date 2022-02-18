package viewmodels

import "github.com/calvine/goauth/core/constants/contact"

type AccountRegisteredTemplateData struct {
	ContactType      contact.Type
	ContactPrincipal string
}
