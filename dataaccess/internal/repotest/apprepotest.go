package repotest

import (
	"context"
	"testing"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

var (
	anotherTestApp models.App

	newTestScope models.Scope
)

const (
	appRepoCreatedByID = "app repos tests"
)

func testAppRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	t.Run("GetAppByID", func(t *testing.T) {
		_testGetAppByID(t, *testHarness.AppRepo)
	})
	t.Run("GetAppsByOwnerID", func(t *testing.T) {
		_testGetAppsByOwnerID(t, *testHarness.AppRepo)
	})
	t.Run("GetAppsByClientID", func(t *testing.T) {
		_testGetAppsByClientID(t, *testHarness.AppRepo)
	})
	t.Run("AddApp", func(t *testing.T) {
		_testAddApp(t, *testHarness.AppRepo)
	})
	t.Run("UpdateApp", func(t *testing.T) {
		_testUpdateApp(t, *testHarness.AppRepo)
	})
	t.Run("DeleteApp", func(t *testing.T) {
		_testDeleteApp(t, *testHarness.AppRepo)
	})

	t.Run("GetScopeByID", func(t *testing.T) {
		_testGetScopeByID(t, *testHarness.AppRepo)
	})
	t.Run("GetScopesByAppID", func(t *testing.T) {
		_testGetScopesByAppID(t, *testHarness.AppRepo)
	})
	t.Run("AddScope", func(t *testing.T) {
		_testAddScope(t, *testHarness.AppRepo)
	})
	t.Run("UpdateScope", func(t *testing.T) {
		_testUpdateScope(t, *testHarness.AppRepo)
	})
	t.Run("DeleteScope", func(t *testing.T) {
		_testDeleteScope(t, *testHarness.AppRepo)
	})
}

func _testGetAppByID(t *testing.T, appRepo repo.AppRepo) {
	app, err := appRepo.GetAppByID(context.TODO(), initialTestApp.ID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get app from underlying data store: %s", err.GetErrorCode())
	}
	if initialTestApp.ID != app.ID {
		t.Errorf("retreived app id does not match expected app is: got %s - expected %s", app.ID, initialTestApp.ID)
	}
}
func _testGetAppsByOwnerID(t *testing.T, appRepo repo.AppRepo) {
	apps, err := appRepo.GetAppsByOwnerID(context.TODO(), initialTestUser.ID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get apps from underlying data store: %s", err.GetErrorCode())
	}
	if len(apps) != 2 {
		t.Errorf("expected to get back two apps based on provided owner id: %v", apps)
	}
}
func _testGetAppsByClientID(t *testing.T, appRepo repo.AppRepo) {
	app, err := appRepo.GetAppByClientID(context.TODO(), initialTestApp.ClientID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get apps from underlying data store: %s", err.GetErrorCode())
	}
	if initialTestApp.ID != app.ID {
		t.Errorf("retreived app id does not match expected app is: got %s - expected %s", app.ID, initialTestApp.ID)
	}
}
func _testAddApp(t *testing.T, appRepo repo.AppRepo) {
	var err errors.RichError
	anotherTestApp, _, err = models.NewApp("fake owner id", "Uber App", "https://uber.app/callback", "https://uber.app/assets/logo")
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create app for test: %s", err.GetErrorCode())
	}
	err = appRepo.AddApp(context.TODO(), &anotherTestApp, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add app to underlying data store: %s", err.GetErrorCode())
	}
	if anotherTestApp.ID == "" {
		t.Error(" app client id should not be empty")
	}
	if anotherTestApp.ClientID == "" {
		t.Error(" app client id should not be empty")
	}
	if anotherTestApp.ClientSecretHash == "" {
		t.Error(" app client secret hash should not be empty")
	}
}
func _testUpdateApp(t *testing.T, appRepo repo.AppRepo) {
	changedAppName := "changed app name"
	anotherTestApp.Name = changedAppName
	err := appRepo.UpdateApp(context.TODO(), &anotherTestApp, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to update app in underlying data store: %s", err.GetErrorCode())
	}
	app, err := appRepo.GetAppByID(context.TODO(), anotherTestApp.ID)
	if err != nil {
		t.Errorf("failed to retreive app from underlying app store for comparison")
	}
	if app.Name != changedAppName {
		t.Errorf("expected app name not correct: got: %s - expected: %s", app.Name, changedAppName)
	}
}
func _testDeleteApp(t *testing.T, appRepo repo.AppRepo) {
	err := appRepo.DeleteApp(context.TODO(), &anotherTestApp, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to delete app from underlying app store: %s", err.GetErrorCode())
	}
}

func _testGetScopeByID(t *testing.T, appRepo repo.AppRepo) {
	scopeID := initialTestAppScopes[0].ID
	scope, err := appRepo.GetScopeByID(context.TODO(), scopeID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get scopes from underlying data store: %s", err.GetErrorCode())
	}
	if scope.ID != scopeID {
		t.Errorf("got scope id that was not expected: got: %s - expected: %s", scope.ID, scopeID)
	}
	if scope.AppID != initialTestApp.ID {
		t.Errorf("got scope app id that was not expected: got: %s - expected: %s", scope.AppID, initialTestApp.ID)
	}
}
func _testGetScopesByAppID(t *testing.T, appRepo repo.AppRepo) {
	scopes, err := appRepo.GetScopesByAppID(context.TODO(), initialTestApp.ID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get scopes from underlying data store: %s", err.GetErrorCode())
	}
	numScopesReturned := len(scopes)
	if numScopesReturned != numScopes {
		t.Errorf("got wrong number of scopes: got %d - expected %d", numScopesReturned, numScopes)
	}
}
func _testAddScope(t *testing.T, appRepo repo.AppRepo) {
	newTestScope = models.NewScope(initialTestApp2.ID, "new_custom_scope", "A scope added for testing add scope and also update scope")
	err := appRepo.AddScope(context.TODO(), &newTestScope, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add scope to underlying data source: %s", err.GetErrorCode())
	}
}
func _testUpdateScope(t *testing.T, appRepo repo.AppRepo) {
	newDescription := "A better description than the previous one"
	newTestScope.Description = newDescription
	err := appRepo.UpdateScope(context.TODO(), &newTestScope, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to update scope in underlying data store: %s", err.GetErrorCode())
	}
}
func _testDeleteScope(t *testing.T, appRepo repo.AppRepo) {
	err := appRepo.DeleteScope(context.TODO(), &newTestScope, appRepoCreatedByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to delete scope from underlying data store: %s", err.GetErrorCode())
	}
}
