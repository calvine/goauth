package repotest

import (
	"testing"

	repo "github.com/calvine/goauth/core/repositories"
)

func testAppRepo(t *testing.T, appRepo repo.AppRepo) {
	t.Run("userRepo.GetAppByID", func(t *testing.T) {
		_testGetAppByID(t, appRepo)
	})
	t.Run("userRepo.GetAppsByOwnerID", func(t *testing.T) {
		_testGetAppsByOwnerID(t, appRepo)
	})
	t.Run("userRepo.AddApp", func(t *testing.T) {
		_testAddApp(t, appRepo)
	})
	t.Run("userRepo.UpdateApp", func(t *testing.T) {
		_testUpdateApp(t, appRepo)
	})
	t.Run("userRepo.DeleteApp", func(t *testing.T) {
		_testDeleteApp(t, appRepo)
	})

	t.Run("userRepo.GetScopesByAppID", func(t *testing.T) {
		_testGetScopesByAppID(t, appRepo)
	})
	t.Run("userRepo.AddScope", func(t *testing.T) {
		_testAddScope(t, appRepo)
	})
	t.Run("userRepo.UpdateScope", func(t *testing.T) {
		_testUpdateScope(t, appRepo)
	})
	t.Run("userRepo.DeleteScope", func(t *testing.T) {
		_testDeleteScope(t, appRepo)
	})
}

func _testGetAppByID(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testGetAppsByOwnerID(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testAddApp(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testUpdateApp(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testDeleteApp(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}

func _testGetScopesByAppID(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testAddScope(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testUpdateScope(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
func _testDeleteScope(t *testing.T, appRepo repo.AppRepo) {
	t.Error("test not implemented")
}
