package models

import (
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
)

type App struct {
	ID               string    `bson:"-"`
	OwnerID          string    `bson:"-"`
	Name             string    `bson:"name"`
	ClientID         string    `bson:"clientId"`
	ClientSecretHash string    `bson:"clientSecret"`
	CallbackURI      string    `bson:"callbackUri"`
	IsDisabled       bool      `bson:"isDisabled"`
	LogoURI          string    `bson:"logoUri"`
	AuditData        auditable `bson:",inline"`
}

func NewApp(ownerID, name, callbackURI, logoURI string) (App, string, errors.RichError) {
	clientID, err := utilities.NewVariableLengthTokenString(1)
	if err != nil {
		return App{}, "", err
	}
	clientSecret, err := utilities.NewVariableLengthTokenString(3)
	if err != nil {
		return App{}, "", err
	}
	// We save the client secret hash so that is not saved in plain text in the database
	// Not using bcrypt because its slow and this needs to be checked per requesat in some cases
	clientSecretHash := utilities.SHA512(clientSecret)
	return App{
		ClientID:         clientID,
		ClientSecretHash: clientSecretHash,
		OwnerID:          ownerID,
		Name:             name,
		CallbackURI:      callbackURI,
		LogoURI:          logoURI,
	}, clientSecret, nil
}
