package memory

import (
	"testing"

	"github.com/calvine/goauth/dataaccess/internal/repotest"
)

func TestMongoRepos(t *testing.T) {
	userRepo := NewMemoryUserRepo()
	contactRepo := NewMemoryContactRepo()
	tokenRepo := NewMemoryTokenRepo()
	testHarnessInput := repotest.RepoTestHarnessInput{
		UserRepo:    &userRepo,
		ContactRepo: &contactRepo,
		TokenRepo:   &tokenRepo,
	}
	repotest.RunReposTestHarness(t, "memory", testHarnessInput)
}
