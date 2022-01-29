package jwt

import (
	"strings"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

type JWTValidator interface {
	GetID() string
	GetSignerFromAlg(alg JWTSigningAlgorithm) (Signer, errors.RichError)
	ValidateHeader(header Header) ([]errors.RichError, bool)
	ValidateClaims(claims StandardClaims) ([]errors.RichError, bool)
	ValidateSignature(algorithm JWTSigningAlgorithm, encodedHeaderAndBody string, signature string) (bool, errors.RichError)
}

type jwtValidator struct {
	id                string
	allowedAlgorithms map[JWTSigningAlgorithm]bool // These are maps to avoid having to loop to find matching items\
	keyIDRequired     bool
	issuerRequired    bool
	allowAnyIssuer    bool
	expectedIssuer    string
	audienceRequired  bool
	allowAnyAudience  bool
	allowedAudience   map[string]bool // These are maps to avoid having to loop to find matching items
	expireRequired    bool
	issuedAtRequired  bool
	notBeforeRequired bool
	subjectRequired   bool
	jtiRequired       bool
	hmacOptions       HMACSigningOptions
	// TODO: add public private key stuff for additional validation types
}

type JWTValidatorOptions struct {
	ID                string
	AllowedAlgorithms []JWTSigningAlgorithm
	KeyIDRequired     bool
	IssuerRequired    bool
	AllowAnyIssuer    bool
	ExpectedIssuer    string
	AudienceRequired  bool
	AllowAnyAudience  bool
	AllowedAudience   []string
	ExpireRequired    bool
	IssuedAtRequired  bool
	NotBeforeRequired bool
	SubjectRequired   bool
	JTIRequired       bool
	HMACOptions       HMACSigningOptions
	// TODO: add public private key stuff for additional validation types
}

// NewJWTValidator creates a JWT validator. I imagine these will end up getting cached if multiple are needed.
func NewJWTValidator(validatorOptions JWTValidatorOptions) (JWTValidator, errors.RichError) {
	if validatorOptions.ID == "" {
		// If we dont get and ID we make one up...
		validatorOptions.ID = uuid.New().String()
	}
	validator := jwtValidator{
		id:                validatorOptions.ID,
		keyIDRequired:     validatorOptions.KeyIDRequired,
		issuerRequired:    validatorOptions.IssuerRequired,
		allowAnyIssuer:    validatorOptions.AllowAnyIssuer,
		expectedIssuer:    validatorOptions.ExpectedIssuer, // should we allow
		expireRequired:    validatorOptions.ExpireRequired,
		audienceRequired:  validatorOptions.AudienceRequired,
		allowAnyAudience:  validatorOptions.AllowAnyAudience,
		issuedAtRequired:  validatorOptions.IssuedAtRequired,
		notBeforeRequired: validatorOptions.NotBeforeRequired,
		subjectRequired:   validatorOptions.SubjectRequired,
		jtiRequired:       validatorOptions.JTIRequired,
		hmacOptions:       validatorOptions.HMACOptions,
	}
	if len(validatorOptions.AllowedAlgorithms) == 0 {
		// You have to specify allowed algorithms
		return validator, coreerrors.NewJWTValidatorNoAlgorithmSpecifiedError(true)
	}
	validator.allowedAlgorithms = make(map[JWTSigningAlgorithm]bool)

	hasMACSecret := len(validator.hmacOptions.Secret) > 0

	for _, a := range validatorOptions.AllowedAlgorithms {
		if !hasMACSecret && strings.HasPrefix(string(a), "HS") {
			return validator, coreerrors.NewJWTValidatorNoHMACSecretProvidedError(true)
		}
		// TODO: have other validation based on the algorithm
		validator.allowedAlgorithms[a] = true
	}

	if !validator.allowAnyAudience {
		if validator.audienceRequired && len(validatorOptions.AllowedAudience) == 0 {
			// you require an audience, but did not allow any audiences
			return validator, coreerrors.NewJWTValidatorAudienceRequiredButNoneProvidedError(true)
		} else {
			validator.allowedAudience = make(map[string]bool)
			for _, a := range validatorOptions.AllowedAudience {
				validator.allowedAudience[a] = true
			}
		}
	}

	if validator.allowAnyIssuer && len(validator.expectedIssuer) != 0 {
		return validator, coreerrors.NewJWTValidatorAllowAnyIssuerAndExpectedIssuerProvidedError(true)
	}

	return validator, nil
}

func (v jwtValidator) GetSignerFromAlg(alg JWTSigningAlgorithm) (Signer, errors.RichError) {
	if strings.HasPrefix(string(alg), "HS") {
		return v.hmacOptions, nil
	}
	// TODO: as other algorithms are added be sure to add here as well
	return nil, coreerrors.NewJWTAlgorithmNotImplementedError(string(alg), true)
}

func (v jwtValidator) GetID() string {
	return v.id
}

func (v jwtValidator) ValidateHeader(header Header) ([]errors.RichError, bool) {
	errs := make([]errors.RichError, 0, 1)

	_, ok := v.allowedAlgorithms[header.Algorithm]
	if !ok {
		errs = append(errs, coreerrors.NewJWTAlgorithmNotAllowedError(string(header.Algorithm), true))
	}

	err := validateKeyID(header.KeyID, v.keyIDRequired)
	if err != nil {
		errs = append(errs, err)
	}

	// not validating type per: https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.9
	// if header.TokenType != Type_JWT {
	// 	// add error
	// }

	return errs, len(errs) == 0
}

func (v jwtValidator) ValidateClaims(claims StandardClaims) ([]errors.RichError, bool) {
	errs := make([]errors.RichError, 0, 4)
	err := validateIssuer(claims.Issuer, v.expectedIssuer, v.issuerRequired, v.allowAnyIssuer)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateExpire(claims.ExpirationTime, v.expireRequired)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateSubject(claims.Subject, v.subjectRequired)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateJti(claims.JWTID, v.jtiRequired)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateIat(claims.IssuedAt, v.issuedAtRequired)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateNbf(claims.NotBefore, v.notBeforeRequired)
	if err != nil {
		errs = append(errs, err)
	}
	err = validateAudience(claims.Audience, v.allowedAudience, v.audienceRequired, v.allowAnyAudience)
	if err != nil {
		errs = append(errs, err)
	}
	return errs, len(errs) == 0
}

func (v jwtValidator) ValidateSignature(alg JWTSigningAlgorithm, encodedHeaderAndBody string, signature string) (bool, errors.RichError) {
	signer, err := v.GetSignerFromAlg(alg)
	if err != nil {
		return false, err
	}
	calculatedSignature, err := signer.Sign(alg, encodedHeaderAndBody)
	if err != nil {
		return false, err
	}
	return signature == calculatedSignature, nil
}

func validateKeyID(keyID string, keyIDRequired bool) errors.RichError {
	if keyIDRequired && len(keyID) == 0 {
		return coreerrors.NewJWTValidatorKeyIDMissingError(true)
	}
	return nil
}

func validateIssuer(issuer, expectedIssuer string, issuerRequired bool, allowAnyIssuer bool) errors.RichError {
	if len(issuer) == 0 {
		if issuerRequired {
			return coreerrors.NewJWTIssuerMissingError(true)
		}
	} else if !allowAnyIssuer && issuer != expectedIssuer {
		return coreerrors.NewJWTIssuerInvalidError(issuer, expectedIssuer, true)
	}
	return nil
}

func validateExpire(exp Time, expireRequired bool) errors.RichError {
	if exp.IsZero() {
		if expireRequired {
			return coreerrors.NewJWTExpireMissingError(true)
		}
	} else if exp.IsInPast() {
		return coreerrors.NewJWTExpiredError(exp.Time(), true)
	}
	return nil
}

func validateSubject(subject string, subjectRequired bool) errors.RichError {
	if subjectRequired && len(subject) == 0 {
		return coreerrors.NewJWTSubjectMissingError(true)
	}
	return nil
}

func validateJti(id string, jtiRequired bool) errors.RichError {
	if jtiRequired && len(id) == 0 {
		return coreerrors.NewJWTIDMissingError(true)
	}
	return nil
}

func validateIat(iat Time, iatRequired bool) errors.RichError {
	if iat.IsZero() {
		if iatRequired {
			return coreerrors.NewJWTIssuedAtMissingError(true)
		}
	} else if iat.IsInFuture() {
		return coreerrors.NewJWTIssuedAtInvalidError(iat.Time(), true)
	}
	return nil
}

func validateNbf(nbf Time, nbfRequired bool) errors.RichError {
	if nbf.IsZero() {
		if nbfRequired {
			return coreerrors.NewJWTNotBeforeMissingError(true)
		}
	} else if nbf.IsInFuture() {
		return coreerrors.NewJWTNotBeforeInFutureError(nbf.Time(), true)
	}
	return nil
}

func validateAudience(audience utilities.CSString, allowedAudiences map[string]bool, audienceRequired, allowAnyAudience bool) errors.RichError {
	if len(audience) == 0 {
		if audienceRequired {
			return coreerrors.NewJWTValidatorAudienceMissingError(true)
		}
	} else if allowAnyAudience {
		return nil
	} else {
		audienceSlice := audience.ToSlice()
		var err errors.RichError
		for _, a := range audienceSlice {
			_, found := allowedAudiences[a]
			if !found {
				// TODO: come back and let this collect all invalid audiences, but for not just one will do...
				err = coreerrors.NewJWTValidatorAudienceInvalidError(a, true)
				break
			}
		}
		return err
	}
	return nil
}
