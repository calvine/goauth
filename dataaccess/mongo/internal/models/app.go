package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreApp models.App

type RepoApp struct {
	ObjectID primitive.ObjectID `bson:"_id"`
	CoreApp  `bson:",inline"`
}

func (ru RepoApp) ToCoreApp() models.App {
	oidString := ru.ObjectID.Hex()
	ru.CoreApp.ID = oidString

	return models.App(ru.CoreApp)
}

func (app CoreApp) ToRepoApp() (RepoApp, errors.RichError) {
	oid, err := primitive.ObjectIDFromHex(app.ID)
	if err != nil {
		return RepoApp{}, coreerrors.NewFailedToParseObjectIDError(app.ID, err, true)
	}
	return RepoApp{
		ObjectID: oid,
		CoreApp:  app,
	}, nil
}

func (app CoreApp) ToRepoAppWithoutID() RepoApp {
	return RepoApp{
		CoreApp: app,
	}
}
