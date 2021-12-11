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

	Alg_NONE = "none" // This should never ever ever ever be used!

	Type_JWT = "JWT"
)

type Header struct {
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
	JWTID          string   `json:"jti,omitempty"` // https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
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
		err := coreerrors.NewJWTSignatureMissingError(encodedJWT, true)
		return nil, err
	}
	return parts, nil
}

func DecodeHeader(encodedHeader string) (Header, errors.RichError) {
	var header Header
	rawHeader, err := Base64UrlDecode(encodedHeader)
	if err != nil {
		errMsg := "failed to base 64 decode header data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedHeader), true)
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

func DecodeBody(encodedBody string) (StandardClaims, errors.RichError) {
	var body StandardClaims
	rawBody, err := Base64UrlDecode(encodedBody)
	if err != nil {
		errMsg := "failed to base 64 decode body data"
		err := coreerrors.NewJWTMalformedError(errMsg, string(encodedBody), true)
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
		// do somthing
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
		// do somthing
		return "", coreerrors.NewJWTEncodingFailedError(err, true)
	}
	encodedClaims := Base64UrlEncode(claimJSONString)
	// removing trailing = signs per https://datatracker.ietf.org/doc/html/rfc7515#section-2 definition of Base64url Encoding
	encodedClaims = strings.TrimRight(encodedClaims, "=")
	return encodedClaims, nil
}

func CalculateHMACSignature(secret string, encodedHeaderAndBody string, hashFunc func() hash.Hash) string {
	hmac := hmac.New(hashFunc, []byte(secret))
	hmac.Write([]byte(encodedHeaderAndBody))
	signatureBytes := hmac.Sum(nil)
	encodedSignature := Base64UrlEncode(signatureBytes)
	return encodedSignature
}

// Base64UrlEncode implemented per https://datatracker.ietf.org/doc/html/rfc7515#appendix-C
func Base64UrlEncode(s []byte) string {
	encodedString := base64.StdEncoding.EncodeToString(s)
	// trim trailing '='
	encodedString = strings.Split(encodedString, "=")[0]
	// convert all '-' to '+'
	encodedString = strings.Replace(encodedString, "+", "-", -1)
	// convert all '/' to '_'
	encodedString = strings.Replace(encodedString, "/", "_", -1)
	return encodedString
}

// Base64UrlDecode implemented per https://datatracker.ietf.org/doc/html/rfc7515#appendix-C
func Base64UrlDecode(encodedString string) ([]byte, errors.RichError) {
	decodedEncodedString := encodedString
	// convert all '+' to '-'
	decodedEncodedString = strings.Replace(decodedEncodedString, "-", "+", -1)
	// convert all '_' to '/'
	decodedEncodedString = strings.Replace(decodedEncodedString, "_", "/", -1)
	// add padding '=' back
	switch len(decodedEncodedString) % 4 {
	case 0:
		// do nothing
	case 2:
		decodedEncodedString += "=="
	case 3:
		decodedEncodedString += "="
	default:
		return nil, coreerrors.NewBase64URLStringInvalidError(encodedString, true)
	}
	decodedString, err := base64.StdEncoding.DecodeString(decodedEncodedString)
	if err != nil {
		return nil, coreerrors.NewBase64DecodeStringFailedError(err, decodedEncodedString, true)
	}
	return []byte(decodedString), nil
}
