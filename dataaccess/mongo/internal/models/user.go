package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreUser models.User

type RepoUser struct {
	ObjectID primitive.ObjectID `bson:"_id"`
	CoreUser `bson:",inline"`
}

func (ru RepoUser) ToCoreUser() models.User {
	oidString := ru.ObjectID.Hex()
	ru.CoreUser.ID = oidString

	return models.User(ru.CoreUser)
}

func (cu CoreUser) ToRepoUser() (RepoUser, errors.RichError) {
	oid, err := primitive.ObjectIDFromHex(cu.ID)
	if err != nil {
		return RepoUser{}, coreerrors.NewFailedToParseObjectIDError(cu.ID, err, true)
	}
	return RepoUser{
		ObjectID: oid,
		CoreUser: cu,
	}, nil
}

func (cu CoreUser) ToRepoUserWithoutID() RepoUser {
	return RepoUser{
		CoreUser: cu,
	}
}
