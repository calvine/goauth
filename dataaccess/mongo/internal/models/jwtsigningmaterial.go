package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreJWTSigningMaterial models.JWTSigningMaterial

type RepoJWTSigningMaterial struct {
	ObjectID               primitive.ObjectID `bson:"_id"`
	CoreJWTSigningMaterial `bson:",inline"`
}

func (ru RepoJWTSigningMaterial) ToCoreJWTSigningMaterial() models.JWTSigningMaterial {
	oidString := ru.ObjectID.Hex()
	ru.CoreJWTSigningMaterial.KeyID = oidString

	return models.JWTSigningMaterial(ru.CoreJWTSigningMaterial)
}

func (cjsm CoreJWTSigningMaterial) ToRepoJWTSigningMaterial() (RepoJWTSigningMaterial, errors.RichError) {
	oid, err := primitive.ObjectIDFromHex(cjsm.KeyID)
	if err != nil {
		return RepoJWTSigningMaterial{}, coreerrors.NewFailedToParseObjectIDError(cjsm.KeyID, err, true)
	}
	return RepoJWTSigningMaterial{
		ObjectID:               oid,
		CoreJWTSigningMaterial: cjsm,
	}, nil
}

func (cjsm CoreJWTSigningMaterial) ToRepoJWTSigningMaterialWithoutID() RepoJWTSigningMaterial {
	return RepoJWTSigningMaterial{
		CoreJWTSigningMaterial: cjsm,
	}
}
