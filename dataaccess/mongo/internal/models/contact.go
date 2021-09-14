package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreContact models.Contact

type RepoContact struct {
	ObjectID    primitive.ObjectID `bson:"id"`
	CoreContact `bson:",inline"`
}

func (rc RepoContact) ToCoreContact() models.Contact {
	oidString := rc.ObjectID.Hex()
	rc.CoreContact.ID = oidString

	return models.Contact(rc.CoreContact)
}

func (cc CoreContact) ToRepoContact() (RepoContact, errors.RichError) {
	oid, err := primitive.ObjectIDFromHex(cc.ID)
	if err != nil {
		return RepoContact{}, coreerrors.NewFailedToParseObjectIDError(cc.ID, err, true)
	}
	return RepoContact{
		ObjectID:    oid,
		CoreContact: cc,
	}, nil
}

func (cc CoreContact) ToRepoContactWithoutID() RepoContact {
	return RepoContact{
		CoreContact: cc,
	}
}
