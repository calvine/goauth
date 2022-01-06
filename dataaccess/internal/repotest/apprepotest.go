package repotest

import (
	"context"
	"testing"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

var (
	testAppOwnerID    string
	testApp           models.App
	testApp2          models.App
	anotherTestApp    models.App
	anotherNewTestApp models.App

	testScope        models.Scope
	anotherTestScope models.Scope
	newTestScope     models.Scope

	nonExistantAppID   string
	nonExistantOwnerID string
	nonExistantScopeID string
)

const (
	appRepoCreatedByID = "app repos tests"
)

func setupAppRepoTestData(t *testing.T, testHarness RepoTestHarnessInput) {
	nonExistantAppID = testHarness.IDGenerator(false)
	nonExistantOwnerID = testHarness.IDGenerator(false)
	nonExistantScopeID = testHarness.IDGenerator(false)
	testAppOwnerID = testHarness.IDGenerator(false)
	_testApp, _, err := models.NewApp(testAppOwnerID, "test app", "uri1", "uri2")
	if err != nil {
		t.Errorf("failed to create app for test: %s", err.Error())
	}
	testApp = _testApp
	_testApp2, _, err := models.NewApp(testAppOwnerID, "test app 2", "uri5", "uri6")
	if err != nil {
		t.Errorf("failed to create app for test: %s", err.Error())
	}
	testApp2 = _testApp2
	_anotherTestApp, _, err := models.NewApp(testHarness.IDGenerator(false), "another test app", "uri3", "uri4")
	if err != nil {
		t.Errorf("failed to create app for test: %s", err.Error())
	}
	anotherTestApp = _anotherTestApp
	_anotherNewTestApp, _, err := models.NewApp(testHarness.IDGenerator(false), "another test app", "uri3", "uri4")
	if err != nil {
		t.Errorf("failed to create app for test: %s", err.Error())
	}
	anotherNewTestApp = _anotherNewTestApp

	// scopes cannot be made here... app ids need to be populated before scopes can be made
}

func setUpScopes() {
	testScope = models.Scope{
		AppID:       testApp.ID,
		Name:        "permission_1",
		DisplayName: "Permission #1",
		Description: "This is permission number 1",
	}
	anotherTestScope = models.Scope{
		AppID:       testApp.ID,
		Name:        "permission_2",
		DisplayName: "Permission #2",
		Description: "This is permission number 2",
	}
	newTestScope = models.Scope{
		AppID:       anotherTestApp.ID,
		Name:        "permission_3",
		DisplayName: "Permission #3",
		Description: "This is permission number 3",
	}
}

func testAppRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	setupAppRepoTestData(t, testHarness)
	t.Run("AddApp", func(t *testing.T) {
		_testAddApp(t, *testHarness.AppRepo)
	})
	// we need to call this so the scop tests will work!
	setUpScopes()
	t.Run("GetAppByID", func(t *testing.T) {
		_testGetAppByID(t, *testHarness.AppRepo)
	})
	t.Run("GetAppsByOwnerID", func(t *testing.T) {
		_testGetAppsByOwnerID(t, *testHarness.AppRepo)
	})
	t.Run("GetAppByClientID", func(t *testing.T) {
		_testGetAppByClientID(t, *testHarness.AppRepo)
	})
	t.Run("UpdateApp", func(t *testing.T) {
		_testUpdateApp(t, *testHarness.AppRepo)
	})
	t.Run("DeleteApp", func(t *testing.T) {
		_testDeleteApp(t, *testHarness.AppRepo)
	})

	t.Run("AddScope", func(t *testing.T) {
		_testAddScope(t, *testHarness.AppRepo)
	})
	t.Run("GetScopeByID", func(t *testing.T) {
		_testGetScopeByID(t, *testHarness.AppRepo)
	})
	t.Run("GetScopesByAppID", func(t *testing.T) {
		_testGetScopesByAppID(t, *testHarness.AppRepo)
	})
	t.Run("UpdateScope", func(t *testing.T) {
		_testUpdateScope(t, *testHarness.AppRepo)
	})
	t.Run("DeleteScope", func(t *testing.T) {
		_testDeleteScope(t, *testHarness.AppRepo)
	})
}

func _testAddApp(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		appToAdd          *models.App
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:     "GIVEN an app to add EXPECT success",
			appToAdd: &testApp,
		},
		{
			name:     "GIVEN another app to add EXPECT success",
			appToAdd: &testApp2,
		},
		{
			name:     "GIVEN one more app to add EXPECT success",
			appToAdd: &anotherTestApp,
		},
		{
			name:     "GIVEN for real this time one more app to add EXPECT success",
			appToAdd: &anotherNewTestApp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := appRepo.AddApp(context.TODO(), tc.appToAdd, appRepoCreatedByID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.appToAdd.ID == "" {
					t.Error("\tapp client id should not be empty")
				}
				if tc.appToAdd.ClientID == "" {
					t.Error("\tapp client id should not be empty")
				}
				if tc.appToAdd.ClientSecretHash == "" {
					t.Error("\tapp client secret hash should not be empty")
				}
			}
		})
	}
}

func _testGetAppByID(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		appID             string
		expectedApp       models.App
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid app id EXPECT success",
			appID:       testApp.ID,
			expectedApp: testApp,
		},
		{
			name:        "GIVEN another valid app id EXPECT success",
			appID:       anotherTestApp.ID,
			expectedApp: anotherTestApp,
		},
		{
			name:              "GIVEN a non existant app id EXPECT error code no app found",
			appID:             nonExistantAppID,
			expectedErrorCode: coreerrors.ErrCodeNoAppFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := appRepo.GetAppByID(context.TODO(), tc.appID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.expectedApp.ID != app.ID {
					t.Errorf("\tretreived app id does not match expected app is: got %s - expected %s", app.ID, tc.expectedApp.ID)
				}
			}
		})
	}
}

func _testGetAppsByOwnerID(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		ownerID           string
		expectedApps      []models.App
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:    "GIVEN a valid owner id with multiple apps EXPECT all apps to be returned",
			ownerID: testAppOwnerID,
			expectedApps: []models.App{
				testApp,
				testApp2,
			},
		},
		{
			name:    "GIVEN a valid owner id with one app EXPECT the proper app to be returned",
			ownerID: anotherTestApp.OwnerID,
			expectedApps: []models.App{
				anotherTestApp,
			},
		},
		{
			name:              "GIVEN an owner id with no associated apps EXPECT error code no apps found",
			ownerID:           nonExistantOwnerID,
			expectedErrorCode: coreerrors.ErrCodeNoAppFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			apps, err := appRepo.GetAppsByOwnerID(context.TODO(), tc.ownerID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numExpected := len(tc.expectedApps)
				numGot := len(apps)
				if numExpected != numGot {
					t.Errorf("\texpected number of apps to get back not expected: got - %d expected - %d", numGot, numExpected)
				}
			}
		})
	}
}

func _testGetAppByClientID(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		clientID          string
		expectedApp       models.App
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid client id EXPECT success",
			clientID:    testApp.ClientID,
			expectedApp: testApp,
		},
		{
			name:              "GIVEN a non existant client id EXPECT error code no app found",
			clientID:          "not a real client id!",
			expectedErrorCode: coreerrors.ErrCodeNoAppFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := appRepo.GetAppByClientID(context.TODO(), tc.clientID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.expectedApp.ID != app.ID {
					t.Errorf("\tretreived app id does not match expected app is: got %s - expected %s", app.ID, tc.expectedApp.ID)
				}
			}
		})
	}
}

func _testUpdateApp(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		appToUpdate       *models.App
		newAppName        string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:        "GIVEN an app to update EXPECT success",
			appToUpdate: &anotherTestApp,
			newAppName:  "a better app name",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.appToUpdate.Name = tc.newAppName
			err := appRepo.UpdateApp(context.TODO(), &anotherTestApp, appRepoCreatedByID)
			if err != nil {
				t.Log(err.Error())
				t.Errorf("\tfailed to update app in underlying data store: %s", err.GetErrorCode())
			}

			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				// pull app back to make sure it updated
				app, rerr := appRepo.GetAppByID(context.TODO(), tc.appToUpdate.ID)
				if rerr != nil {
					t.Errorf("\tfailed to retreive app from underlying app store for comparison: %s", rerr.Error())
				}
				if app.Name != tc.newAppName {
					t.Errorf("\texpected app name not correct: got: %s - expected: %s", app.Name, tc.newAppName)
				}
			}
		})
	}
}

func _testDeleteApp(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		appToDelete       *models.App
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:        "GIVEN an app to delete EXPECT success",
			appToDelete: &anotherTestApp,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := appRepo.DeleteApp(context.TODO(), tc.appToDelete, appRepoCreatedByID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} //else {
			// non failure test code here
			//}
		})
	}
}

func _testAddScope(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		scopeToAdd        *models.Scope
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:       "GIVEN a scope to add EXPECT success",
			scopeToAdd: &testScope,
		},
		{
			name:       "GIVEN another scope to add EXPECT success",
			scopeToAdd: &anotherTestScope,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := appRepo.AddScope(context.TODO(), tc.scopeToAdd, appRepoCreatedByID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} //else {
			// non failure test code here
			//}
		})
	}
}

func _testGetScopeByID(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		scopeID           string
		expectedScope     models.Scope
		expectedAppID     string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:          "GIVEN EXPECT success",
			scopeID:       testScope.ID,
			expectedScope: testScope,
			expectedAppID: testApp.ID,
		},
		{
			name:              "GIVEN a non existant scope id EXPECT error code no scope found",
			scopeID:           nonExistantScopeID,
			expectedErrorCode: coreerrors.ErrCodeNoScopeFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scope, err := appRepo.GetScopeByID(context.TODO(), tc.scopeID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if scope.ID != tc.expectedScope.ID {
					t.Errorf("\tgot scope id that was not expected: got: %s - expected: %s", scope.ID, tc.expectedScope.ID)
				}
				if scope.AppID != tc.expectedAppID {
					t.Errorf("\tgot scope app id that was not expected: got: %s - expected: %s", scope.AppID, tc.expectedAppID)
				}
			}
		})
	}
}

func _testGetScopesByAppID(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		appID             string
		expectedScopes    []models.Scope
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:  "GIVEN an app id with scopes EXPECT all scopes to be returned",
			appID: testApp.ID,
			expectedScopes: []models.Scope{
				testScope,
				anotherTestScope,
			},
		},
		{
			name:              "GIVEN an app id with no scopes EXPECT error code no scope found",
			appID:             anotherNewTestApp.ID,
			expectedErrorCode: coreerrors.ErrCodeNoScopeFound,
		},
		{
			name:              "GIVEN a non existant app id EXPECT error code no scope found",
			appID:             nonExistantAppID,
			expectedErrorCode: coreerrors.ErrCodeNoScopeFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scopes, err := appRepo.GetScopesByAppID(context.TODO(), tc.appID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numExpectedScopes := len(tc.expectedScopes)
				numScopesReturned := len(scopes)
				if numScopesReturned != numExpectedScopes {
					t.Errorf("\tgot wrong number of scopes: got %d - expected %d", numScopesReturned, numExpectedScopes)
				}
			}
		})
	}
}

func _testUpdateScope(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name                string
		scopeToUpdate       *models.Scope
		newScopeName        string
		newScopeDisplayName string
		newScopeDescription string
		expectedErrorCode   string
	}
	testCases := []testCase{
		{
			name:                "GIVEN EXPECT success",
			scopeToUpdate:       &testScope,
			newScopeName:        "updated_scope_name",
			newScopeDisplayName: "Updated Scope Name!",
			newScopeDescription: "This is an updated scope",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.scopeToUpdate.Name = tc.newScopeName
			tc.scopeToUpdate.DisplayName = tc.newScopeDisplayName
			tc.scopeToUpdate.Description = tc.newScopeDescription
			err := appRepo.UpdateScope(context.TODO(), tc.scopeToUpdate, appRepoCreatedByID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				scope, rerr := appRepo.GetScopeByID(context.TODO(), tc.scopeToUpdate.ID)
				if rerr != nil {
					t.Errorf("\tfailed to retreive app from underlying app store for comparison: %s", rerr.Error())
				}
				if scope.Name != tc.newScopeName {
					t.Errorf("\tname value was not expected: got - %s expected - %s", scope.Name, tc.newScopeName)
				}
				if scope.DisplayName != tc.newScopeDisplayName {
					t.Errorf("\tdisplay name value was not expected: got - %s expected - %s", scope.DisplayName, tc.newScopeDisplayName)
				}
				if scope.Description != tc.newScopeDescription {
					t.Errorf("\tdescription value was not expected: got - %s expected - %s", scope.Description, tc.newScopeDescription)
				}
			}
		})
	}
}

func _testDeleteScope(t *testing.T, appRepo repo.AppRepo) {
	type testCase struct {
		name              string
		scopeToDelete     *models.Scope
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:          "GIVEN a scope to delete EXPECT success",
			scopeToDelete: &testScope,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := appRepo.DeleteScope(context.TODO(), tc.scopeToDelete, appRepoCreatedByID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} //else {
			// non failure test code here
			//}
		})
	}
}
