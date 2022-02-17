package core

type AuthStatus int

const (
	Unauthenticated AuthStatus = iota
	Invalid
	Expired
	Authenticated
)

type TemplateType string

const (
	Text TemplateType = "text"
	Html TemplateType = "html"
)

type ContactType string

//TODO: come back to things like this and pull out into aliased type or enum
const (
	Email  ContactType = "email"
	Mobile ContactType = "mobile"
)
