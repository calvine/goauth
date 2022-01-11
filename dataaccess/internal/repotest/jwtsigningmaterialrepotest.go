package repotest

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

const (
	jwtSigningMaterialRepoCreatedByID = "memory jwt signing material repo test"
)

var (
	nonExistantJWTSigningMaterialID string

	jwtSigningMaterial1 models.JWTSigningMaterial
	jwtSigningMaterial2 models.JWTSigningMaterial
	jwtSigningMaterial3 models.JWTSigningMaterial
	jwtSigningMaterial4 models.JWTSigningMaterial
	jwtSigningMaterial5 models.JWTSigningMaterial
)

func setupJWTSigningMaterialRepoTestData(_ *testing.T, testingHarness RepoTestHarnessInput) {
	nonExistantJWTSigningMaterialID = testingHarness.IDGenerator(false)
	jwtSigningMaterial1 = models.JWTSigningMaterial{
		AlgorithmType: "HMAC",
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret",
		},
		Expiration: nullable.NullableTime{
			HasValue: false,
		},
		Disabled: false,
	}
	jwtSigningMaterial2 = models.JWTSigningMaterial{
		AlgorithmType: "HMAC",
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret2",
		},
		Expiration: nullable.NullableTime{
			HasValue: true,
			Value:    time.Now().UTC().Add(time.Hour),
		},
		Disabled: false,
	}
	jwtSigningMaterial3 = models.JWTSigningMaterial{
		AlgorithmType: "OTHER",
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret3",
		},
		Expiration: nullable.NullableTime{
			HasValue: true,
			Value:    time.Now().Add(time.Hour),
		},
		Disabled: true,
	}
	jwtSigningMaterial4 = models.JWTSigningMaterial{
		AlgorithmType: "HMAC",
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret2",
		},
		Expiration: nullable.NullableTime{
			HasValue: true,
			Value:    time.Now().Add(-time.Hour),
		},
		Disabled: false,
	}
	jwtSigningMaterial5 = models.JWTSigningMaterial{
		AlgorithmType: "HMAC",
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    "testsecret2",
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
	t.Run("GetValidJWTSigningMaterialByAlgorithmType", func(t *testing.T) {
		_testGetValidJWTSigningMaterialByAlgorithmType(t, *testHarness.JWTSigningMaterialRepo)
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
		{
			name:                 "GIVEN disabled jwt signing material EXPECT success",
			signingMaterialToAdd: &jwtSigningMaterial4,
		},
		{
			name:                 "GIVEN disabled jwt signing material EXPECT success",
			signingMaterialToAdd: &jwtSigningMaterial5,
		},
		// {
		// 	name: "GIVEN jwt signing material with a duplicate key id EXPECT error code jwt signing material key id not unique",
		// 	signingMaterialToAdd: &models.JWTSigningMaterial{
		// 		KeyID: jwtSigningMateriaKeyID1, // this was added in the first test case
		// 	},
		// 	expectedErrorCode: coreerrors.ErrCodeJWTSigningMaterialKeyIDNotUnique,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := jwtSigningMaterialRepo.AddJWTSigningMaterial(context.TODO(), tc.signingMaterialToAdd, jwtSigningMaterialRepoCreatedByID)
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
	testCases := []testCase{
		{
			name:                       "GIVEN given a valid key id for a jwt signing material EXPECT success",
			keyID:                      jwtSigningMaterial1.KeyID,
			expectedJWTSigningMaterial: jwtSigningMaterial1,
		},
		{
			name:              "GIVEN given an invalid key id for a jwt signing material EXPECT error code no jwt signing material found",
			keyID:             nonExistantJWTSigningMaterialID,
			expectedErrorCode: coreerrors.ErrCodeNoJWTSigningMaterialFound,
		},
	}
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
				if jsm.HMACSecret.Value != tc.expectedJWTSigningMaterial.HMACSecret.Value {
					t.Errorf("\tsecret is not expected value: got - %s expected - %s", jsm.HMACSecret.Value, tc.expectedJWTSigningMaterial.HMACSecret.Value)
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
}

func _testGetValidJWTSigningMaterialByAlgorithmType(t *testing.T, jwtSigningMaterialRepo repo.JWTSigningMaterialRepo) {
	type testCase struct {
		name                       string
		algorithmType              string
		expectedJWTSigningMaterial []models.JWTSigningMaterial
		expectedErrorCode          string
	}
	testCases := []testCase{
		{
			name:          "GIVEN given a valid algorithm type for a jwt signing material EXPECT success",
			algorithmType: "HMAC",
			expectedJWTSigningMaterial: []models.JWTSigningMaterial{
				jwtSigningMaterial1,
				jwtSigningMaterial2,
			},
		},
		{
			name:                       "GIVEN given a non existant algorithm type for a jwt signing material EXPECT no jwt signing material to be returned",
			algorithmType:              "zzz",
			expectedJWTSigningMaterial: []models.JWTSigningMaterial{},
			//expectedErrorCode: coreerrors.ErrCodeNoJWTSigningMaterialFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsms, err := jwtSigningMaterialRepo.GetValidJWTSigningMaterialByAlgorithmType(context.TODO(), tc.algorithmType)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numGot := len(jsms)
				numExpected := len(tc.expectedJWTSigningMaterial)
				if numGot != numExpected {
					t.Errorf("\tnumber of jwt signing material returned not expected: got - %d expected - %d", numGot, numExpected)
				}
				for _, ejsm := range tc.expectedJWTSigningMaterial {
					found := false
					for _, jsm := range jsms {
						if jsm.KeyID == ejsm.KeyID {
							found = true
							if ejsm.KeyID != jsm.KeyID {
								t.Errorf("\tkey id is not expected value: got - %s expected - %s", jsm.KeyID, ejsm.KeyID)
							}
							if ejsm.HMACSecret.Value != jsm.HMACSecret.Value {
								t.Errorf("\tsecret is not expected value: got - %s expected - %s", jsm.HMACSecret.Value, ejsm.HMACSecret.Value)
							}
							if ejsm.Expiration.Value.Round(time.Second) != jsm.Expiration.Value.Round(time.Second) {
								t.Errorf("\texpiration is not expected value: got - %v expected - %v", jsm.Expiration.Value, ejsm.Expiration.Value)
							}
							if ejsm.Disabled != jsm.Disabled {
								t.Errorf("\tdisabled is not expected value: got - %v expected - %v", jsm.Disabled, ejsm.Disabled)
							}
							break
						}
					}
					if !found {
						t.Errorf("\tfailed to find jwt signing material: %+v", ejsm)
						continue
					}
				}
			}
		})
	}
}
