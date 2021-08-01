package utilities

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func InterleaveStrings(s1, s2 string) string {
	// TODO: do somthing more interesting than concatenating, to make reverse engineering more difficult.
	return s1 + s2
}

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
