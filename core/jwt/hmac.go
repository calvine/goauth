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

func (hso HMACSigningOptions) Sign(alg JWTSigningAlgorithm, encodedHeaderAndBody string) (string, errors.RichError) {
	// hmac := hmac.New(hashFunc, []byte(secret))
	var hashFunc HashFunc
	switch alg {
	case HS256:
		hashFunc = sha256.New
	case HS384:
		hashFunc = sha512.New384
	case HS512:
		hashFunc = sha512.New
	default:
		return "", coreerrors.NewJWTAlgorithmNotAllowedError(string(alg), true)
	}
	hmac := hmac.New(hashFunc, []byte(hso.Secret))
	hmac.Write([]byte(encodedHeaderAndBody))
	signatureBytes := hmac.Sum(nil)
	encodedSignature := Base64UrlEncode(signatureBytes)
	return encodedSignature, nil
}

func (hso HMACSigningOptions) GetAlgorithmFamily() JWTSingingAlgorithmFamily {
	return HMAC
}

func (hso HMACSigningOptions) IsAlgorithmSupported(alg JWTSigningAlgorithm) bool {
	return alg == HS256 || alg == HS384 || alg == HS512
}
