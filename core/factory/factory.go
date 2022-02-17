package factory

import (
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/richerror/errors"
)

// TODO: think about how to implement this interface?
type JWTFactory interface {
	NewSignedJWT(sub string, aud []string, exp nullable.NullableDuration) (jwt.JWT, string, errors.RichError)
	NewSignedJWTWithSigner(sub string, aud []string, exp nullable.NullableDuration, signer jwt.Signer) (jwt.JWT, string, errors.RichError)
}

// JWTValidatorFactory is used to create validatators for jwts.
type JWTValidatorFactory interface {
	NewJWTValidator(keyID string) (jwt.JWTValidator, errors.RichError)
	NewJWTValidatorWithSigner(keyID string, signer jwt.Signer) (jwt.JWTValidator, errors.RichError)
	NewValidatorWithOptions(options jwt.JWTValidatorOptions) (jwt.JWTValidator, errors.RichError)
}

type ServiceLinkFactory interface {
	CreateLink(linkPath string, queryParams map[string]string) (string, errors.RichError)
	CreatePasswordResetLink(passwordResetToken string) (string, errors.RichError)
	CreateConfirmContactLink(confirmContactToken string) (string, errors.RichError)
	CreateLoginLink() (string, errors.RichError)
	CreateMagicLoginLink(magicLoginToken string) (string, errors.RichError)
	CreateUserRegisterLink() (string, errors.RichError)
	// Eventually app management urls
}
