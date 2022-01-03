package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

var (
	nonExistantJWTSigningMaterialID string

	jwtSigningMaterial1 models.JWTSigningMaterial
	jwtSigningMaterial2 models.JWTSigningMaterial
	jwtSigningMaterial3 models.JWTSigningMaterial
)

func setupJWTSigningMaterialRepoTestData(_ *testing.T, testingHarness RepoTestHarnessInput) {
	nonExistantJWTSigningMaterialID = testingHarness.IDGenerator(false)
	jwtSigningMaterial1 = models.JWTSigningMaterial{
		KeyID: "123",
		Secret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret",
		},
		Expiration: nullable.NullableTime{
			HasValue: false,
		},
		Disabled: false,
	}
	jwtSigningMaterial2 = models.JWTSigningMaterial{
		KeyID: "456",
		Secret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret2",
		},
		Expiration: nullable.NullableTime{
			HasValue: true,
			Value:    time.Now().Add(time.Hour),
		},
		Disabled: false,
	}
	jwtSigningMaterial3 = models.JWTSigningMaterial{
		KeyID: "789",
		Secret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret3",
		},
		Expiration: nullable.NullableTime{
			HasValue: true,
			Value:    time.Now().Add(time.Hour),
		},
		Disabled: true,
	}
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
	testCases := []testCase{
		{
			name:                 "GIVEN jwt signing material with no expiration EXPECT success",
			signingMaterialToAdd: &jwtSigningMaterial1,
		},
		{
			name:                 "GIVEN jwt signing material EXPECT success",
			signingMaterialToAdd: &jwtSigningMaterial2,
		},
		{
			name:                 "GIVEN disabled jwt signing material EXPECT success",
			signingMaterialToAdd: &jwtSigningMaterial3,
		},
	}
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
