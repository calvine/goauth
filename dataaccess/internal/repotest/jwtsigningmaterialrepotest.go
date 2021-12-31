package repotest

import (
	"context"
	"testing"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

var (
	nonExistantJWTSigningMaterialID string

	jwtSigningMaterial1 models.JWTSigningMaterial = models.JWTSigningMaterial{}
	jwtSigningMaterial2 models.JWTSigningMaterial = models.JWTSigningMaterial{}
	jwtSigningMaterial3 models.JWTSigningMaterial = models.JWTSigningMaterial{}
)

func setupJWTSigningMaterialRepoTestData(_ *testing.T, testingHarness RepoTestHarnessInput) {
	nonExistantJWTSigningMaterialID = testingHarness.IDGenerator(false)
}

func testJWTSigningMaterialRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	setupJWTSigningMaterialRepoTestData(t, testHarness)
	t.Run("AddJWTSigningMaterial", func(t *testing.T) {
		_testAddJWTSigningMaterial(t, *testHarness.JWTSigningMaterialRepo)
	})
	t.Run("GetJWTSigningMaterialByKeyID", func(t *testing.T) {
		_testGetJWTSigningMaterialByKeyID(t, *testHarness.JWTSigningMaterialRepo)
	})
}

func _testAddJWTSigningMaterial(t *testing.T, jwtSigningMaterialRepo repo.JWTSigningMaterialRepo) {
	type testCase struct {
		name                 string
		signingMaterialToAdd *models.JWTSigningMaterial
		expectedErrorCode    string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := jwtSigningMaterialRepo.AddJWTSigningMaterial(context.TODO(), tc.signingMaterialToAdd, "")
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.signingMaterialToAdd.AuditData.CreatedByID == "" {
					t.Error("\texpected created by id to be populated")
				}
				if tc.signingMaterialToAdd.AuditData.CreatedOnDate.IsZero() {
					t.Error("\texpected created on date to be populated")
				}
			}
		})
	}
	t.Error("test not implemented")
}

func _testGetJWTSigningMaterialByKeyID(t *testing.T, jwtSigningMaterialRepo repo.JWTSigningMaterialRepo) {
	type testCase struct {
		name                       string
		keyID                      string
		expectedJWTSigningMaterial models.JWTSigningMaterial
		expectedErrorCode          string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsm, err := jwtSigningMaterialRepo.GetJWTSigningMaterialByKeyID(context.TODO(), tc.keyID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if jsm.KeyID != tc.expectedJWTSigningMaterial.KeyID {
					t.Errorf("\tkey id is not expected value: got - %s expected - %s", jsm.KeyID, tc.expectedJWTSigningMaterial.KeyID)
				}
				if jsm.Secret.Value != tc.expectedJWTSigningMaterial.Secret.Value {
					t.Errorf("\tsecret is not expected value: got - %s expected - %s", jsm.Secret.Value, tc.expectedJWTSigningMaterial.Secret.Value)
				}
				if jsm.Expiration.Value != tc.expectedJWTSigningMaterial.Expiration.Value {
					t.Errorf("\texpiration is not expected value: got - %v expected - %v", jsm.Expiration.Value, tc.expectedJWTSigningMaterial.Expiration.Value)
				}
				if jsm.Disabled != tc.expectedJWTSigningMaterial.Disabled {
					t.Errorf("\tdisabled is not expected value: got - %v expected - %v", jsm.Disabled, tc.expectedJWTSigningMaterial.Disabled)
				}
			}
		})
	}
	t.Error("test not implemented")
}
