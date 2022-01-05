package models

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreScope models.Scope

type RepoScope struct {
	ObjectID  primitive.ObjectID `bson:"_id"`
	CoreScope `bson:",inline"`
}

func (ru RepoScope) ToCoreScope() models.Scope {
	oidString := ru.ObjectID.Hex()
	ru.CoreScope.ID = oidString

	return models.Scope(ru.CoreScope)
}

func (scope CoreScope) ToRepoScope() (RepoScope, errors.RichError) {
	oid, err := primitive.ObjectIDFromHex(scope.ID)
	if err != nil {
		return RepoScope{}, coreerrors.NewFailedToParseObjectIDError(scope.ID, err, true)
	}
	return RepoScope{
		ObjectID:  oid,
		CoreScope: scope,
	}, nil
}

func (scope CoreScope) ToRepoScopeWithoutID() RepoScope {
	return RepoScope{
		CoreScope: scope,
	}
}
