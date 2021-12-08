package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"hash"
	"strings"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

// JSON Web Token (JWT) 		https://datatracker.ietf.org/doc/html/rfc7519
// JSON Web Signature (JWS) 	https://datatracker.ietf.org/doc/html/rfc7515
// JSON Web Encryption (JWE)	https://datatracker.ietf.org/doc/html/rfc7516
// JSON Web Algorithms (JWA)	https://datatracker.ietf.org/doc/html/rfc7518

const (
	Alg_HS256 = "HS256"
	Alg_HS384 = "HS384"
	Alg_HS512 = "HS512"
	// TODO: implement RS, ES and PS based algorithms

	Type_JWT = "JWT"
)

type HeaderFields struct {
	Algorithm   string `json:"alg"` // https://datatracker.ietf.org/doc/html/rfc7518#section-3.1
	ContentType string `json:"cty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-5.2
	TokenType   string `json:"typ"` // TODO: use this... https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.9
}

type StandardClaims struct {
	Issuer         string   `json:"iss"`           // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
	Subject        string   `json:"sub"`           // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
	Audience       []string `json:"aud,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	ExpirationTime Time     `json:"exp,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
	NotBefore      Time     `json:"nbf,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
	IssuedAt       Time     `json:"iat,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
	JWTID          string   `json:"jwi,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
}

func splitEncodedJWT(encodedJWT string) ([]string, errors.RichError) {
	parts := strings.Split(encodedJWT, ".")
	if len(parts) != 3 {
		errMsg := "there should be exactly three parts"
		err := coreerrors.NewJWTMalformedError(errMsg, encodedJWT, true)
		return nil, err
	}
	if len(parts[2]) == 0 {
		// no signature provided
		err := coreerrors.NewJWTMissingSignatureError(encodedJWT, true)
		return nil, err
	}
	return parts, nil
}

func DecodeHeader(encodedHeader string) (HeaderFields, errors.RichError) {
	var header HeaderFields
	rawHeader, err := base64.StdEncoding.DecodeString(encodedHeader)
	if err != nil {
		errMsg := "failed to base 64 decode header data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedHeader), true)
		return header, err
	}
	err = json.Unmarshal(rawHeader, &header)
	if err != nil {
		errMsg := "failed to unmarshal header data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedHeader), true)
		return header, err
	}
	return header, nil
}

func DecodeBody(encodedBody string) (StandardClaims, errors.RichError) {
	var body StandardClaims
	rawBody, err := base64.StdEncoding.DecodeString(encodedBody)
	if err != nil {
		errMsg := "failed to base 64 decode body data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedBody), true)
		return body, err
	}
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		errMsg := "failed to unmarshal body data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(rawBody), true)
		return body, err
	}
	return body, nil
}

func (header HeaderFields) Encode() (string, errors.RichError) {
	headerJSONString, err := json.Marshal(header)
	if err != nil {
		// do somthing
		return "", coreerrors.NewJWTEncodingFailedError(err, true)
	}
	encodedClaims := base64.StdEncoding.EncodeToString(headerJSONString)
	return encodedClaims, nil
}

func (claims StandardClaims) Encode() (string, errors.RichError) {
	claimJSONString, err := json.Marshal(claims)
	if err != nil {
		// do somthing
		return "", coreerrors.NewJWTEncodingFailedError(err, true)
	}
	encodedClaims := base64.StdEncoding.EncodeToString(claimJSONString)
	return encodedClaims, nil
}

func calculateHMACSignature(secret string, encodedHeaderAndBody string, hashFunc func() hash.Hash) string {
	hmac := hmac.New(hashFunc, []byte(secret))
	signatureBytes := hmac.Sum([]byte(encodedHeaderAndBody))
	encodedSignature := base64.StdEncoding.EncodeToString(signatureBytes)
	return encodedSignature
}
