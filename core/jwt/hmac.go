package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

type HMACSigningOptions struct {
	Secret string
}

func NewHMACSigningOptions(s string) (HMACSigningOptions, errors.RichError) {
	if len(s) == 0 {
		return HMACSigningOptions{}, coreerrors.NewHMACSigingOptionsMissingSecretError(true)
	}
	return HMACSigningOptions{
		Secret: s,
	}, nil
}

func (hso HMACSigningOptions) Sign(alg string, encodedHeaderAndBody string) (string, errors.RichError) {
	// hmac := hmac.New(hashFunc, []byte(secret))
	var hashFunc HashFunc
	switch alg {
	case Alg_HS256:
		hashFunc = sha256.New
	case Alg_HS384:
		hashFunc = sha512.New384
	case Alg_HS512:
		hashFunc = sha512.New
	default:
		return "", coreerrors.NewJWTAlgorithmNotAllowedError(alg, true)
	}
	hmac := hmac.New(hashFunc, []byte(hso.Secret))
	hmac.Write([]byte(encodedHeaderAndBody))
	signatureBytes := hmac.Sum(nil)
	encodedSignature := Base64UrlEncode(signatureBytes)
	return encodedSignature, nil
}
