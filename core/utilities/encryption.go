package utilities

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func SHA256(input string) (string, error) {
	hash := sha256.Sum256([]byte(input))
	hashString := hex.EncodeToString(hash[:])
	return string(hashString), nil
}

func SHA512(input string) (string, error) {
	hash := sha512.Sum512([]byte(input))
	hashString := hex.EncodeToString(hash[:])
	return string(hashString), nil
}
