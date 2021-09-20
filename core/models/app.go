package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
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

func ValidateApp(includeID bool, app App) errors.RichError {
	fields := make(map[string]interface{})
	if includeID && app.ID == "" {
		fields["ID"] = "app ID cannot be empty"
	}
	if app.ClientID == "" {
		fields["ClientID"] = "app ClientID cannot be empty"
	}
	if app.ClientSecretHash == "" {
		fields["ClientSecretHash"] = "app ClientSecretHash cannot be empty"
	}
	if app.Name == "" {
		fields["Name"] = "app Name cannot be empty"
	}
	if app.CallbackURI == "" {
		fields["CallbackURI"] = "app CallbackURI cannot be empty"
	}
	if app.LogoURI == "" {
		fields["LogoURI"] = "app LogoURI cannot be empty"
	}

	if len(fields) > 0 {
		return coreerrors.NewInvalidAppCreationError(fields, false)
	}
	return nil
}
