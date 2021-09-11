package models

type App struct {
	ID           string `bson:"-"`
	OwnerID      string `bson:"-"`
	ClientID     string `bson:"clientId"`
	ClientSecret string `bson:"clientSecret"`
	CallbackURI  string `bson:"callbackUri"`
	IsDisabled   bool   `bson:"isDisabled"`
	LogoURI      string `bson:"logoUri"`
}
