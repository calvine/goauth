package utilities

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
	"golang.org/x/crypto/bcrypt"
)

func BcryptHashString(s string, cost int) (string, errors.RichError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), cost)
	if err != nil {
		return "", coreerrors.NewBcryptPasswordHashErrorError("", err, true)
	}
	return string(hash), nil
}

func BcryptCompareStringAndHash(hash, s, assetID string) (bool, errors.RichError) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else if err != nil {
		return false, coreerrors.NewBcryptPasswordHashErrorError(assetID, err, true)
	}
	return true, nil
}
