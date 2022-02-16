package jwt

import (
	"encoding/json"
	"hash"
	"strings"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

// JSON Web Token (JWT) 		https://datatracker.ietf.org/doc/html/rfc7519
// JSON Web Signature (JWS) 	https://datatracker.ietf.org/doc/html/rfc7515
// JSON Web Encryption (JWE)	https://datatracker.ietf.org/doc/html/rfc7516
// JSON Web Algorithms (JWA)	https://datatracker.ietf.org/doc/html/rfc7518

type HashFunc func() hash.Hash

type Header struct {
	Algorithm   JWTSigningAlgorithm `json:"alg"`           // https://datatracker.ietf.org/doc/html/rfc7518#section-3.1
	TokenType   string              `json:"typ"`           // TODO: use this... https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.9
	ContentType string              `json:"cty,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-5.2
	KeyID       string              `json:"kid,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.4`
}

type StandardClaims struct {
	Issuer         string             `json:"iss,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Subject        string             `json:"sub,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Audience       utilities.CSString `json:"aud,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	ExpirationTime Time               `json:"exp,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	NotBefore      Time               `json:"nbf,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	IssuedAt       Time               `json:"iat,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	JWTID          string             `json:"jti,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
	Scopes         utilities.CSString `json:"scopes,omitempty"`
}

// Signer is an interface for somthing that can sign a JWT.
// TODO: we need somthing to verify signature. for HMAC stuff thats the same as signing, but for pub/pri key pair stuff it will be different...
type Signer interface {
	// Sign produces a signature for a given encoded header and body with the given algorithm
	Sign(alg JWTSigningAlgorithm, encodedHeaderAndBody string) (string, errors.RichError)
	// Verify TODO: implement this?
	// Verify(token JWT) (string, errors.RichError)
	// GetAlgorithmFamily returns the algorithm family the signer belongs to
	GetAlgorithmFamily() JWTSigningAlgorithmFamily
	// IsAlgorithmSupported returns a boolean value indicating if the algorithm provided is supported
	IsAlgorithmSupported(alg JWTSigningAlgorithm) bool
}

// func (s Signer) MarshalJSON() ([]byte, error) {
// 	switch t := s.(type) {
// 	case HMACSigningOptions:
// 		return json.Marshal(t)
// 	default:
// 		return nil, coreerrors.NewJWTAlgorithmNotImplementedError(string(s.GetAlgorithmFamily()), true)
// 	}
// }

// func (s *Signer) UnmarshalJSON(data []byte, v interface{}) error {
// 	var tempMap map[string]*json.RawMessage
// 	err := json.Unmarshal(data, &tempMap)
// 	if err != nil {
// 		// TODO: make a rich error for failing to unmarshal...
// 		return err
// 	}
// 	t, ok := tempMap["type"]
// 	if !ok {
// 		// FIXME: feels weird not to use code gen, but this is a very specific error that I dont think will be used elsewhere. still go back and make this a generated error.
// 		return errors.NewRichError("SignerUnmarshalMissingTypeField", "json data did not have a type, must not be a signer").WithStack(1)
// 	}
// 	switch t {
// 	case HMAC:

// 	default:
// 		return coreerrors.NewJWTAlgorithmNotImplementedError(t, true)
// 	}
// }

type JWT struct {
	Header    Header
	Claims    StandardClaims
	Signature string
}

// JWTSigningAlgorithm are specific individual jwt signing algorithms
type JWTSigningAlgorithm string

// JWTSigningAlgorithmFamily defines a family of algorithms that contains multiple signing algorithms
type JWTSigningAlgorithmFamily string

const (
	HMAC JWTSigningAlgorithmFamily = "HMAC"

	HS256 JWTSigningAlgorithm = "HS256"
	HS384 JWTSigningAlgorithm = "HS384"
	HS512 JWTSigningAlgorithm = "HS512"

	// TODO: implement RS, ES and PS based algorithms

	NONE JWTSigningAlgorithm = "none" // This should never ever ever ever be used!

	Typ_JWT = "JWT"
)

var (
	HMACAlgorithms = []JWTSigningAlgorithm{HS256, HS384, HS512}
)

func NewUnsignedJWT(alg JWTSigningAlgorithm, iss string, aud []string, sub string, duration nullable.NullableDuration, notBefore time.Time) JWT {
	header := Header{
		Algorithm: alg,
		TokenType: Typ_JWT,
	}

	exp := Time{}
	if duration.HasValue {
		exp = FromDuration(duration.Value)
	}

	claims := StandardClaims{
		Issuer:         iss,
		Audience:       utilities.NewCSString(aud), // if no audience is provided this will be an empty string and opitted from the token for json marshaling
		IssuedAt:       NewTime(),
		Subject:        sub,
		ExpirationTime: exp,
		NotBefore:      Time(notBefore),
		JWTID:          uuid.NewString(),
	}

	return JWT{
		Header: header,
		Claims: claims,
	}
}

func DecodeAndValidateJWT(encodedJWT string, validator JWTValidator) (JWT, errors.RichError) {
	token, err := Decode(encodedJWT)
	if err != nil {
		return JWT{}, err
	}

	// validate signature
	valid, err := validator.ValidateSignature(token.Header.Algorithm, encodedJWT)
	if err != nil {
		return JWT{}, err
	}
	if !valid {
		return JWT{}, coreerrors.NewJWTSignatureInvalidError(encodedJWT, nil, true)
	}
	// validate claims
	claimErrors, valid := validator.ValidateClaims(token.Claims)
	if !valid {
		rErr := coreerrors.NewJWTStandardClaimsInvalidError(encodedJWT, true)
		for _, e := range claimErrors {
			rErr.AddError(e)
		}
		return JWT{}, rErr
	}
	return token, nil

}

// ENcodeSignedJWT encodes a signed JWT. If no signature is set an error is returned
func (jwt JWT) EncodeSignedJWT() (string, errors.RichError) {
	parts := make([]string, 0, 3)
	encodedHeader, err := jwt.Header.Encode()
	if err != nil {
		return "", err
	}
	encodedBody, err := jwt.Claims.Encode()
	if err != nil {
		return "", err
	}
	parts = append(parts, encodedHeader, encodedBody)
	if len(jwt.Signature) == 0 {
		return "", coreerrors.NewJWTSignatureMissingError(strings.Join(parts, "."), true)
	}
	parts = append(parts, jwt.Signature)
	return strings.Join(parts, "."), nil
}

func (jwt *JWT) SignAndEncode(signer Signer) (string, errors.RichError) {
	parts := make([]string, 0, 3)
	encodedHeader, err := jwt.Header.Encode()
	if err != nil {
		return "", err
	}
	encodedBody, err := jwt.Claims.Encode()
	if err != nil {
		return "", err
	}
	parts = append(parts, encodedHeader, encodedBody)
	encodedHeaderAndBody := strings.Join(parts[:2], ".")
	// NOTE: this code is not intended work with the none algorithm. it should not be used...
	signature, err := signer.Sign(jwt.Header.Algorithm, encodedHeaderAndBody)
	if err != nil {
		return "", err
	}
	jwt.Signature = signature
	parts = append(parts, signature)
	return strings.Join(parts, "."), nil
}

func SplitEncodedJWT(encodedJWT string) ([]string, errors.RichError) {
	parts := strings.Split(encodedJWT, ".")
	if len(parts) != 3 {
		errMsg := "there should be exactly three parts"
		err := coreerrors.NewJWTMalformedError(errMsg, encodedJWT, true)
		return nil, err
	}
	// FIXME: I dont think this can be reached...
	if len(parts[2]) == 0 {
		// no signature provided
		err := coreerrors.NewJWTSignatureMissingError(encodedJWT, true)
		return nil, err
	}
	return parts, nil
}

func DecodeJWTPartRaw(part string) ([]byte, errors.RichError) {
	raw, err := Base64UrlDecode(part)
	if err != nil {
		errMsg := "failed to base 64 decode data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(part), true)
		return nil, err
	}
	return raw, nil
}

func Decode(encodedJWT string) (JWT, errors.RichError) {
	var token JWT
	parts, err := SplitEncodedJWT(encodedJWT)
	if err != nil {
		return JWT{}, err
	}
	token.Header, err = DecodeHeader(parts[0])
	if err != nil {
		return JWT{}, err
	}
	token.Claims, err = DecodeStandardClaims(parts[1])
	if err != nil {
		return JWT{}, err
	}
	token.Signature = parts[2]
	return token, nil
}

func DecodeHeader(encodedHeader string) (Header, errors.RichError) {
	var header Header
	rawHeader, err := DecodeJWTPartRaw(encodedHeader)
	if err != nil {
		return header, err
	}
	nerr := json.Unmarshal(rawHeader, &header)
	if nerr != nil {
		errMsg := "failed to unmarshal header data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedHeader), true)
		return header, err
	}
	return header, nil
}

func DecodeStandardClaims(encodedBody string) (StandardClaims, errors.RichError) {
	var body StandardClaims
	rawBody, err := DecodeJWTPartRaw(encodedBody)
	if err != nil {
		return body, err
	}
	nerr := json.Unmarshal(rawBody, &body)
	if nerr != nil {
		errMsg := "failed to unmarshal body data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(rawBody), true)
		return body, err
	}
	return body, nil
}

func (header Header) Encode() (string, errors.RichError) {
	headerJSONString, err := json.Marshal(header)
	if err != nil {
		return "", coreerrors.NewJWTEncodingFailedError(err, true)
	}
	encodedClaims := Base64UrlEncode(headerJSONString)
	// removing trailing = signs per https://datatracker.ietf.org/doc/html/rfc7515#section-2 definition of Base64url Encoding
	encodedClaims = strings.TrimRight(encodedClaims, "=")
	return encodedClaims, nil
}

func (claims StandardClaims) Encode() (string, errors.RichError) {
	claimJSONString, err := json.Marshal(claims)
	if err != nil {
		return "", coreerrors.NewJWTEncodingFailedError(err, true)
	}
	encodedClaims := Base64UrlEncode(claimJSONString)
	// removing trailing = signs per https://datatracker.ietf.org/doc/html/rfc7515#section-2 definition of Base64url Encoding
	encodedClaims = strings.TrimRight(encodedClaims, "=")
	return encodedClaims, nil
}
