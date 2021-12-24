package core

type AuthStatus int

const (
	Unauthenticated AuthStatus = iota
	Invalid
	Expired
	Authenticated
)

//TODO: come back to things like this and pull out into aliased type or enum
const (
	CONTACT_TYPE_EMAIL  = "email"
	CONTACT_TYPE_MOBILE = "mobile"
)
