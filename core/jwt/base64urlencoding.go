package jwt

import (
	"encoding/base64"
	"strings"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

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
