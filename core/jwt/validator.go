package jwt

import (
	"crypto/sha256"
	"crypto/sha512"
	"strings"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

type JWTValidator interface {
	ValidateHeader(header HeaderFields) ([]errors.RichError, bool)
	ValidateClaims(claims StandardClaims) ([]errors.RichError, bool)
	ValidateSignature(alg string, encodedHeaderAndBody string, signature string) (bool, errors.RichError)
}

type jwtValidator struct {
	id                string
	allowedAlgorithms map[string]bool // These are maps to avoid having to loop to find matching items
	expectedIssuer    string
	allowedAudience   map[string]bool // These are maps to avoid having to loop to find matching items
	audienceRequired  bool
	expireRequired    bool
	issuedAtRequired  bool
	notBeforeRequired bool
	subjectRequired   bool
	jwiRequired       bool
	hmacSecret        string
	// TODO: add public private key stuff for additional validation
}

type JWTValidatorOptions struct {
	ID                string
	AllowedAlgorithms []string
	ExpectedIssuer    string
	AllowedAudience   []string
	AudienceRequired  bool
	ExpireRequired    bool
	IssuedAtRequired  bool
	NotBeforeRequired bool
	SubjectRequired   bool
	JWIRequired       bool
	HMACSecret        string
	// TODO: add public private key stuff for additional validation
}

// NewJWTValidator creates a JWT validator. I imagine these will end up getting cached if multiple are needed.
func NewJWTValidator(validatorOptions JWTValidatorOptions) (JWTValidator, errors.RichError) {
	validator := jwtValidator{
		id:                validatorOptions.ID,
		expectedIssuer:    validatorOptions.ExpectedIssuer,
		expireRequired:    validatorOptions.ExpireRequired,
		audienceRequired:  validatorOptions.AudienceRequired,
		issuedAtRequired:  validatorOptions.IssuedAtRequired,
		notBeforeRequired: validatorOptions.NotBeforeRequired,
		subjectRequired:   validatorOptions.SubjectRequired,
		jwiRequired:       validatorOptions.JWIRequired,
		hmacSecret:        validatorOptions.HMACSecret,
	}
	if len(validator.allowedAlgorithms) == 0 {
		// You have to specify allowed algorithms
		return validator, coreerrors.NewJWTValidatorNoAlgorithmSpecifiedError(true)
	}
	validator.allowedAlgorithms = make(map[string]bool)

	validatHMACSecret := len(validator.hmacSecret) > 0

	for _, a := range validatorOptions.AllowedAlgorithms {
		if !validatHMACSecret && strings.HasPrefix(a, "HS") {
			return validator, coreerrors.NewJWTValidatorNoHMACSecretProvidedError(true)
		}
		// TODO have other validation based on
		validator.allowedAlgorithms[a] = true
	}

	if validator.audienceRequired {
		validator.allowedAudience = make(map[string]bool)
		for _, a := range validatorOptions.AllowedAudience {
			validator.allowedAlgorithms[a] = true
		}
	}

	return validator, nil
}

func (v jwtValidator) ValidateHeader(header HeaderFields) ([]errors.RichError, bool) {
	errs := make([]errors.RichError, 0, 2)
	valid := true

	_, ok := v.allowedAlgorithms[header.Algorithm]
	if !ok {
		errs = append(errs, coreerrors.NewJWTAlgorithmNotAllowedError(header.Algorithm, true))
	}

	// not validating type per: https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.9
	// if header.TokenType != Type_JWT {
	// 	// add error
	// }

	return errs, valid
}

func (v jwtValidator) ValidateClaims(claims StandardClaims) ([]errors.RichError, bool) {
	errs := make([]errors.RichError, 0, 4)
	valid := true
	if len(claims.Subject) == 0 {
		// subject must be populated
		err := coreerrors.NewJWTMissingSubjectError(nil, true)
		errs = append(errs, err)
		valid = false
	}
	now := time.Now()
	exp := claims.ExpirationTime.Time()
	if !exp.IsZero() && exp.Before(now) {
		// token is expired
		err := coreerrors.NewJWTExipredError(exp, nil, true)
		errs = append(errs, err)
		valid = false
	}
	iat := claims.IssuedAt.Time()
	if iat.After(now) {
		// issued at is in the future some how...
		err := coreerrors.NewJWTInvalidIssuedAtError(iat, nil, true)
		errs = append(errs, err)
		valid = false
	}
	nbf := claims.NotBefore.Time()
	if nbf.Before(now) {
		// token not before has not yet passed
		err := coreerrors.NewJWTNotBeforeInFutureError(nbf, nil, true)
		errs = append(errs, err)
		valid = false
	}
	// TODO: validate audience, issuer, and that jwt id is populated?Its a
	return errs, valid
}

func (v jwtValidator) ValidateSignature(alg string, encodedHeaderAndBody string, signature string) (bool, errors.RichError) {
	var calculatedSignature string
	switch alg {
	case Alg_HS256:
		calculatedSignature = CalculateHMACSignature(v.hmacSecret, encodedHeaderAndBody, sha256.New)
	case Alg_HS384:
		calculatedSignature = CalculateHMACSignature(v.hmacSecret, encodedHeaderAndBody, sha512.New384)
	case Alg_HS512:
		calculatedSignature = CalculateHMACSignature(v.hmacSecret, encodedHeaderAndBody, sha512.New)
	default:
		return false, coreerrors.NewJWTAlgorithmNotImplementedError(alg, true)
	}
	return signature == calculatedSignature, nil
}
