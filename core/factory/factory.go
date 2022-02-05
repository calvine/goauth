package factory

import (
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/richerror/errors"
)

type JWTValidatorFactory interface {
	NewJWTValidator(keyID string) (jwt.JWTValidator, errors.RichError)
	NewValidatorWithOptions(options jwt.JWTValidatorOptions) (jwt.JWTValidator, errors.RichError)
}
