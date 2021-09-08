package models

import "time"

// go:generate stringer -type=TokenType
type TokenType int

const (
	TokenTypeInvalid TokenType = iota
	TokenTypeCSRF
	TokenTypeConfirmContact
	TokenTypePasswordReset
)

// Token is a temporary item that can be used as a shared secret like a password reset token or a confirm contact token. They can be tide to a target entity like a user to ensure they are consumed by the proper targets.
type Token struct {
	// Value needs to a be a universially unique value like a uuid or something like that. This is the token passed around.
	Value string
	// TokenType is the type of token the token is.
	TokenType TokenType
	// Expiration is the time at which the token expires. all tokens must expire, so this must have a value.
	Expiration time.Time
	// TargetID is to specify the entity to who the token applies. if the token can be accessed anonymously, leave this blank.
	TargetID string
	// MetaData is a map that contains general purpose data related to a token.
	MetaData map[string]string
}

func NewToken(value, targetID string, tokenType TokenType, validFor time.Duration) Token {
	return Token{
		Value:      value,
		TargetID:   targetID,
		TokenType:  tokenType,
		Expiration: time.Now().Add(validFor),
	}
}

func (t Token) WithMetaData(metaData map[string]string) Token {
	t.MetaData = metaData
	return t
}

func (t Token) AddMetaData(key, value string) Token {
	if t.MetaData == nil {
		t.MetaData = make(map[string]string, 0)
	}
	t.MetaData[key] = value
	return t
}
