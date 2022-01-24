package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
	"github.com/google/uuid"
)

const (
	ALGTYP_HMAC = "HMAC"
)

type JWTSigningMaterial struct {
	ID            string                  `bson:"-"`
	KeyID         string                  `bson:"keyId"`
	AlgorithmType string                  `bson:"algorithmType"`
	HMACSecret    nullable.NullableString `bson:"hmacSecret"`
	Expiration    nullable.NullableTime   `bson:"expiration"`
	Disabled      bool                    `bson:"disabled"`
	AuditData     auditable               `bson:",inline"`
	// PublicKey nullable.NullableString `bson:"publicKey"`
	// PrivateKey nullable.NullableString `bson:"privateKey"`
}

func NewHMACJWTSigningMaterial(secret string, expiration nullable.NullableTime) JWTSigningMaterial {
	return JWTSigningMaterial{
		KeyID:         uuid.Must(uuid.NewRandom()).String(), // TODO: make a function to create random unique key ids?
		AlgorithmType: ALGTYP_HMAC,
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    secret,
		},
		Expiration: expiration,
		Disabled:   false,
	}
}

func (jsm *JWTSigningMaterial) IsExpired() bool {
	now := time.Now().UTC()
	if jsm.Expiration.HasValue && jsm.Expiration.Value.Before(now) {
		return true
	}
	return false
}
