package models

import (
	"time"

	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
)

//go:generate stringer -type=TokenType -output=token_type_string.go
type TokenType int

const (
	TokenTypeInvalid TokenType = iota
	TokenTypeCSRF
	TokenTypeConfirmContact
	TokenTypePasswordReset
	TokenTypeSession
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

func NewToken(targetID string, tokenType TokenType, validFor time.Duration) (Token, errors.RichError) {
	var token Token
	value, err := utilities.NewTokenString()
	if err != nil {
		return token, err
	}
	return Token{
		Value:      value,
		TargetID:   targetID,
		TokenType:  tokenType,
		Expiration: time.Now().Add(validFor),
	}, nil
}

func (t *Token) WithMetaData(metaData map[string]string) {
	t.MetaData = metaData
}

func (t *Token) AddMetaData(key, value string) {
	if t.MetaData == nil {
		t.MetaData = make(map[string]string, 0)
	}
	t.MetaData[key] = value
}

func (t Token) IsExpired() bool {
	return t.Expiration.Before(time.Now())
}
