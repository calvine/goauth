package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

type JWTSigningMaterial struct {
	KeyID         string                  `bson:"-"`
	AlgorithmType string                  `bson:"algorithmType"`
	HMACSecret    nullable.NullableString `bson:"hmacSecret"`
	Expiration    nullable.NullableTime   `bson:"expiration"`
	Disabled      bool                    `bson:"disabled"`
	AuditData     auditable               `bson:",inline"`
	// PublicKey nullable.NullableString `bson:"publicKey"`
	// PrivateKey nullable.NullableString `bson:"privateKey"`
}

func (jsm *JWTSigningMaterial) IsExpired() bool {
	now := time.Now().UTC()
	if jsm.Expiration.HasValue && jsm.Expiration.Value.Before(now) {
		return true
	}
	return false
}
