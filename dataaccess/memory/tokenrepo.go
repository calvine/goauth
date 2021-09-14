package memory

import (
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

type localTokenRepo struct {
	tokenMap map[string]models.Token
}

func NewMemoryTokenRepo() repo.TokenRepo {
	tokenMap := make(map[string]models.Token)
	return &localTokenRepo{tokenMap}
}

func (ltr *localTokenRepo) GetToken(tokenValue string) (models.Token, errors.RichError) {
	token, ok := ltr.tokenMap[tokenValue]
	if !ok {
		return token, coreerrors.NewTokenNotFoundError(tokenValue, true)
	}
	return token, nil
}

func (ltr *localTokenRepo) PutToken(token models.Token) errors.RichError {
	ltr.tokenMap[token.Value] = token
	return nil
}

func (ltr *localTokenRepo) DeleteToken(tokenValue string) errors.RichError {
	delete(ltr.tokenMap, tokenValue)
	return nil
}
