package memory

import (
	"testing"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/dataaccess/internal/repotest"
	"github.com/google/uuid"
)

func TestMemoryRepos(t *testing.T) {
	users := make(map[string]models.User)
	contacts := make(map[string]models.Contact)
	userRepo, err := NewMemoryUserRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	contactRepo, err := NewMemoryContactRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appRepo := NewMemoryAppRepo()
	tokenRepo := NewMemoryTokenRepo()
	testHarnessInput := repotest.RepoTestHarnessInput{
		UserRepo:    &userRepo,
		ContactRepo: &contactRepo,
		AppRepo:     &appRepo,
		TokenRepo:   &tokenRepo,
		IDGenerator: func(getZeroId bool) string {
			if getZeroId {
				return uuid.UUID{}.String()
			}
			return uuid.Must(uuid.NewRandom()).String()
		},
	}
	repotest.RunReposTestHarness(t, "memory", testHarnessInput)
}
