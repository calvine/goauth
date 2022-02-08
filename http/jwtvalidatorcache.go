package http

import (
	"time"

	"github.com/calvine/goauth/core/jwt"
)

type jwtValidatorCache map[string]cachedJWTValidator

type cachedJWTValidator struct {
	validator  jwt.JWTValidator
	expiration time.Time
}

func (jvc *jwtValidatorCache) GetCachedValidator(id string) (jwt.JWTValidator, bool) {
	v, ok := (*jvc)[id]
	if !ok {
		return nil, false
	}
	now := time.Now()
	if v.expiration.After(now) {
		return nil, false
	}
	return v.validator, true
}

func (jvc *jwtValidatorCache) CacheValidator(id string, validator jwt.JWTValidator, cacheDuration time.Duration) {
	// for now we dont care if the cached validator already exists we will just overwrite
	(*jvc)[id] = cachedJWTValidator{
		expiration: time.Now().Add(cacheDuration),
		validator:  validator,
	}
}
