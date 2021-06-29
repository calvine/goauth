package models

import (
	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreContact models.Contact

type RepoContact struct {
	ObjectId    primitive.ObjectID `bson:"id"`
	CoreContact `bson:",inline"`
}

func (rc RepoContact) ToCoreContact() models.Contact {
	oidString := rc.ObjectId.Hex()
	rc.CoreContact.Id = oidString

	return models.Contact(rc.CoreContact)
}

func (cc CoreContact) ToRepoContact() (RepoContact, error) {
	oid, err := primitive.ObjectIDFromHex(cc.Id)
	if err != nil {
		return RepoContact{}, err
	}
	return RepoContact{
		ObjectId:    oid,
		CoreContact: cc,
	}, nil
}

func (cc CoreContact) ToRepoContactWithoutId() RepoContact {
	return RepoContact{
		CoreContact: cc,
	}
}
