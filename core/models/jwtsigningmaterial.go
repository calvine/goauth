package models

import "github.com/calvine/goauth/core/nullable"

type JWTSigningMaterial struct {
	KeyID      string                  `bson:"keyID"`
	Secret     nullable.NullableString `bson:"secret"`
	Expiration nullable.NullableTime   `bson:"expiration"`
	Disabled   bool                    `bson:"disabled"`
	AuditData  auditable               `bson:",inline"`
	// PublicKey nullable.NullableString `bson:"publicKey"`
	// PrivateKey nullable.NullableString `bson:"privateKey"`
}
