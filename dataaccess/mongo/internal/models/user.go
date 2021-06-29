package models

import (
	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreUser models.User

type RepoUser struct {
	ObjectId primitive.ObjectID `bson:"_id"`
	CoreUser `bson:",inline"`
}

func (ru RepoUser) ToCoreUser() models.User {
	oidString := ru.ObjectId.Hex()
	ru.CoreUser.Id = oidString

	return models.User(ru.CoreUser)
}

func (cu CoreUser) ToRepoUser() (RepoUser, error) {
	oid, err := primitive.ObjectIDFromHex(cu.Id)
	if err != nil {
		return RepoUser{}, err
	}
	return RepoUser{
		ObjectId: oid,
		CoreUser: cu,
	}, nil
}

func (cu CoreUser) ToRepoUserWithoutId() RepoUser {
	return RepoUser{
		CoreUser: cu,
	}
}
