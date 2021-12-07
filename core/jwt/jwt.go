package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

const (
	Alg_HS256 = "HS256"
	Alg_HS512 = "HS512"
	// TODO: implement ES and PS based algorithms

	Type_JWT = "JWT"
)

type HeaderFields struct {
	Algorithm string `json:"alg"`
	TokenType string `json:"typ"`
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

func (claims StandardClaims) ValidateClaims() errors.RichError {
	if len(claims.Subject) == 0 {
		// subject must be populated
	}
	now := time.Now()
	exp := claims.ExpirationTime.Time()
	if exp.Before(now) {
		// token is expired
	}
	iat := claims.IssuedAt.Time()
	if iat.After(now) {
		// issued at is in the future some how...
	}
	nbf := claims.NotBefore.Time()
	if nbf.Before(now) {
		// token not before has not yet passed
	}
	// TODO: validate audience, issuer, and that jwt id is populated?
	return nil
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

func decodeHeader(encodedHeader string) (HeaderFields, errors.RichError) {
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

func decodeBody(encodedBody string) (StandardClaims, errors.RichError) {
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
