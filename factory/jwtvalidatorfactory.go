package factory

import (
	corefactory "github.com/calvine/goauth/core/factory"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/richerror/errors"
)

type jwtValidatorFactory struct {
	defaultOptions jwt.JWTValidatorOptions
}

func NewJWTValidatorFactory(defaultOptions jwt.JWTValidatorOptions) corefactory.JWTValidatorFactory {
	return jwtValidatorFactory{
		defaultOptions: defaultOptions,
	}
}

func (jvf jwtValidatorFactory) NewJWTValidator(keyID string) (jwt.JWTValidator, errors.RichError) {
	options := jvf.defaultOptions
	return jwt.NewJWTValidator(options)
}

func (jvf jwtValidatorFactory) NewValidatorWithOptions(options jwt.JWTValidatorOptions) (jwt.JWTValidator, errors.RichError) {
	return jwt.NewJWTValidator(options)
}
