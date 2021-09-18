package utilities

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"strings"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

func InterleaveStrings(s1, s2 string) string {
	// TODO: do somthing more interesting than concatenating, to make reverse engineering more difficult.
	return s1 + s2
}

func SHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])
	return string(hashString)
}

func SHA512(input string) string {
	hash := sha512.Sum512([]byte(input))
	hashString := hex.EncodeToString(hash[:])
	return string(hashString)
}

func NewTokenString() (string, errors.RichError) {
	var tokenString string
	for i := 0; i < 2; i++ {
		uuid, err := uuid.NewRandom()
		tokenString += strings.ReplaceAll(uuid.String(), "-", "")
		if err != nil {
			return "", coreerrors.NewGenerateUUIDFailedError(err, true)
		}
	}
	return tokenString, nil
}

func NewVariableLengthTokenString(numGuids int) (string, errors.RichError) {
	var tokenString string
	if numGuids <= 0 {
		numGuids = 2
	}
	for i := 0; i < numGuids; i++ {
		uuid, err := uuid.NewRandom()
		tokenString += strings.ReplaceAll(uuid.String(), "-", "")
		if err != nil {
			return "", coreerrors.NewGenerateUUIDFailedError(err, true)
		}
	}
	return tokenString, nil
}
