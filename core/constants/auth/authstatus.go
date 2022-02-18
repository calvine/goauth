package auth

type Status int

const (
	Unauthenticated Status = iota
	Invalid
	Expired
	Authenticated
)
