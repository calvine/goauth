package factory

import (
	coreerrors "github.com/calvine/goauth/core/errors"
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
	options.ID = keyID
	return jwt.NewJWTValidator(options)
}

func (jvf jwtValidatorFactory) NewJWTValidatorWithSigner(keyID string, signer jwt.Signer) (jwt.JWTValidator, errors.RichError) {
	options := jvf.defaultOptions
	options.RemoveSignerOptions()
	// I am not a fan of this, see comments on the jwt.JWTValidatorOptions struct
	switch o := signer.(type) {
	case jwt.HMACSigningOptions:
		options.HMACOptions = o
	default:
		// TODO: as other signer options as added make sure to add the code to clear them here.
		return nil, coreerrors.NewJWTAlgorithmNotAllowedError(string(signer.GetAlgorithmFamily()), true)
	}
	options.ID = keyID
	validator, err := jwt.NewJWTValidator(options)
	if err != nil {
		return nil, err
	}
	return validator, nil
}

func (jvf jwtValidatorFactory) NewValidatorWithOptions(options jwt.JWTValidatorOptions) (jwt.JWTValidator, errors.RichError) {
	return jwt.NewJWTValidator(options)
}
