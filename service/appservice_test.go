package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
	"github.com/calvine/richerror/errors"
)

const (
	createdByAppService = "app service tests"
	numScopesToMake     = 10
	testAppOneOwnerID   = "testapponeownerid"
	testAppTwoOwnerID   = "testapptwoownerid"
)

var (
	testAppOne_One             models.App
	testAppOne_OneScopes       []models.Scope
	testAppOne_OneClientSecret string

	testAppOne_Two             models.App
	testAppOne_TwoScopes       []models.Scope
	testAppOne_TwoClientSecret string

	testAppOne_Three             models.App
	testAppOne_ThreeScopes       []models.Scope
	testAppOne_ThreeClientSecret string

	testAppTwo models.App
	// testAppTwoScopes       []models.Scope
	testAppTwoClientSecret string
)

func setupAppServiceTestData(t *testing.T, appRepo repo.AppRepo) {
	var err errors.RichError
	testAppOne_One, testAppOne_OneClientSecret, err = models.NewApp(testAppOneOwnerID, "", "", "")
	if err != nil {
		t.Errorf("failed to create test app one: %s", err.Error())
	}
	err = appRepo.AddApp(context.TODO(), &testAppOne_One, createdByAppService)
	if err != nil {
		t.Errorf("failed to add app to underlying data store: %s", err.Error())
	}
	testAppOne_OneScopes = make([]models.Scope, 0, numScopesToMake)
	for i := 1; i <= numScopesToMake; i++ {
		scope := models.NewScope(testAppOne_One.ID, fmt.Sprintf("test_app_one_one_scope_%d", i), fmt.Sprintf("test app one_one scope %d", i))
		err := appRepo.AddScope(context.TODO(), &scope, createdByAppService)
		if err != nil {
			t.Errorf("failed to add scope to app with id %s: %s", testAppOne_One.ID, err.Error())
		}
		testAppOne_OneScopes = append(testAppOne_OneScopes, scope)
	}

	testAppOne_Two, testAppOne_TwoClientSecret, err = models.NewApp(testAppOneOwnerID, "", "", "")
	if err != nil {
		t.Errorf("failed to create test app one: %s", err.Error())
	}
	err = appRepo.AddApp(context.TODO(), &testAppOne_Two, createdByAppService)
	if err != nil {
		t.Errorf("failed to add app to underlying data store: %s", err.Error())
	}
	testAppOne_TwoScopes = make([]models.Scope, 0, numScopesToMake)
	for i := 1; i <= numScopesToMake; i++ {
		scope := models.NewScope(testAppOne_Two.ID, fmt.Sprintf("test_app_one_two_scope_%d", i), fmt.Sprintf("test app one_two scope %d", i))
		err := appRepo.AddScope(context.TODO(), &scope, createdByAppService)
		if err != nil {
			t.Errorf("failed to add scope to app with id %s: %s", testAppOne_Two.ID, err.Error())
		}
		testAppOne_TwoScopes = append(testAppOne_TwoScopes, scope)
	}

	testAppOne_Three, testAppOne_ThreeClientSecret, err = models.NewApp(testAppOneOwnerID, "", "", "")
	if err != nil {
		t.Errorf("failed to create test app one: %s", err.Error())
	}
	err = appRepo.AddApp(context.TODO(), &testAppOne_Three, createdByAppService)
	if err != nil {
		t.Errorf("failed to add app to underlying data store: %s", err.Error())
	}
	testAppOne_ThreeScopes = make([]models.Scope, 0, numScopesToMake)
	for i := 1; i <= numScopesToMake; i++ {
		scope := models.NewScope(testAppOne_Three.ID, fmt.Sprintf("test_app_one_three_scope_%d", i), fmt.Sprintf("test app one_three scope %d", i))
		err := appRepo.AddScope(context.TODO(), &scope, createdByAppService)
		if err != nil {
			t.Errorf("failed to add scope to app with id %s: %s", testAppOne_Three.ID, err.Error())
		}
		testAppOne_ThreeScopes = append(testAppOne_ThreeScopes, scope)
	}

	testAppTwo, testAppTwoClientSecret, err = models.NewApp(testAppOneOwnerID, "", "", "")
	if err != nil {
		t.Errorf("failed to create test app one: %s", err.Error())
	}
	err = appRepo.AddApp(context.TODO(), &testAppTwo, createdByAppService)
	if err != nil {
		t.Errorf("failed to add app to underlying data store: %s", err.Error())
	}
}

func buildAppService(t *testing.T) services.AppService {
	appRepo := memory.NewMemoryAppRepo()
	auditLogRepo := memory.NewMemoryAuditLogRepo(false)
	appService := NewAppService(appRepo, auditLogRepo)
	setupAppServiceTestData(t, appRepo)
	return appService
}

func TestAppService(t *testing.T) {
	appService := buildAppService(t)
	t.Run("GetAppsByOwnerID", func(t *testing.T) {
		_testGetAppsByOwnerID(t, appService)
	})
	t.Run("GetAppByID", func(t *testing.T) {
		_testGetAppByID(t, appService)
	})
	t.Run("GetAppByClientID", func(t *testing.T) {
		_testGetAppByClientID(t, appService)
	})
	t.Run("GetAppAndScopesByClientID", func(t *testing.T) {
		_testGetAppAndScopesByClientID(t, appService)
	})
	t.Run("AddApp", func(t *testing.T) {
		_testAddApp(t, appService)
	})
	t.Run("UpdateApp", func(t *testing.T) {
		_testUpdateApp(t, appService)
	})
	t.Run("DeleteApp", func(t *testing.T) {
		_testDeleteApp(t, appService)
	})
	t.Run("GetScopeByID", func(t *testing.T) {
		_testGetScopeByID(t, appService)
	})
	t.Run("GetScopesByAppID", func(t *testing.T) {
		_testGetScopesByAppID(t, appService)
	})
	t.Run("GetScopesByClientID", func(t *testing.T) {
		_testGetScopesByClientID(t, appService)
	})
	t.Run("AddScopesToApp", func(t *testing.T) {
		_testAddScopesToApp(t, appService)
	})
	t.Run("UpdateScope", func(t *testing.T) {
		_testUpdateScope(t, appService)
	})
	t.Run("DeleteScope", func(t *testing.T) {
		_testDeleteScope(t, appService)
	})
}

func _testGetAppsByOwnerID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no apps found
}

func _testGetAppByID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no app found
}

func _testGetAppByClientID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no client id found
}

func _testGetAppAndScopesByClientID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure not client id found
}

func _testAddApp(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no name

	// failure no owner id

	// failure no callback uri

	// failure no logo uri
}

func _testUpdateApp(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no name

	// failure no owner id

	// failure no callback uri

	// failure no logo uri
}

func _testDeleteApp(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure not app found
}

func _testGetScopeByID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	//failure no scope id found
}

func _testGetScopesByAppID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no app id found
}

func _testGetScopesByClientID(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no client id found
}

func _testAddScopesToApp(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no app found

	// failure no app id

	// failure no name

	// failure no description
}

func _testUpdateScope(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure no app id

	// failure no name

	// failure no description
}

func _testDeleteScope(t *testing.T, appService services.AppService) {
	t.Error("test not implemented")
	// success

	// failure scope found
}
