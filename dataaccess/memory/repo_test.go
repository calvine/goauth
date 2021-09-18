package memory

import (
	"testing"

	"github.com/calvine/goauth/dataaccess/internal/repotest"
)

func TestMemoryRepos(t *testing.T) {
	userRepo := NewMemoryUserRepo()
	contactRepo := NewMemoryContactRepo()
	appRepo := NewMemoryAppRepo()
	tokenRepo := NewMemoryTokenRepo()
	testHarnessInput := repotest.RepoTestHarnessInput{
		UserRepo:    &userRepo,
		ContactRepo: &contactRepo,
		AppRepo:     &appRepo,
		TokenRepo:   &tokenRepo,
	}
	repotest.RunReposTestHarness(t, "memory", testHarnessInput)
}
