package service

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	coreservices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
	"github.com/calvine/goauth/internal/testutils"
	"github.com/google/uuid"
	"go.uber.org/zap/zaptest"
)

const (
	JSMSTestCreatedBy = "jwtsigningmaterialservicecreatedby"
)

var (
	jwtSigningMaterial1         models.JWTSigningMaterial
	jwtSigningMaterial2         models.JWTSigningMaterial
	jwtSigningMaterial3         models.JWTSigningMaterial
	duplicateJWTSigningMaterial models.JWTSigningMaterial
)

func TestJWTSigningMaterialService(t *testing.T) {
	jsms := buildJWTSigningMaterialService(t)
	t.Run("_testAddJWTSigningMaterial", func(t *testing.T) {
		_testAddJWTSigningMaterial(t, jsms)
	})
	t.Run("_testGetJWTSigningMaterialByKeyID", func(t *testing.T) {
		_testGetJWTSigningMaterialByKeyID(t, jsms)
	})
	t.Run("_testGetValidJWTSigningMaterialByAlgorithmType", func(t *testing.T) {
		_testGetValidJWTSigningMaterialByAlgorithmType(t, jsms)
	})
}

func buildJWTSigningMaterialService(t *testing.T) coreservices.JWTSigningMaterialService {
	jsmr := memory.NewMemoryJWTSigningMaterialRepo()
	jsms := NewJWTSigningMaterialService(jsmr)
	setupJWTSigningMaterialServiceTestData(t)
	return jsms
}

func setupJWTSigningMaterialServiceTestData(t *testing.T) {
	jwtSigningMaterial1 = models.NewHMACJWTSigningMaterial("test1", nullable.NullableTime{
		HasValue: false,
		Value:    time.Time{},
	})
	jwtSigningMaterial2 = models.NewHMACJWTSigningMaterial("test2", nullable.NullableTime{
		HasValue: false,
		Value:    time.Time{},
	})
	jwtSigningMaterial3 = models.JWTSigningMaterial{
		KeyID:         "other key",
		AlgorithmType: "OTHER",
		Expiration:    nullable.NullableTime{},
		Disabled:      false,
	}
	duplicateJWTSigningMaterial = models.JWTSigningMaterial{
		KeyID:         "other key",
		AlgorithmType: "OTHER",
		Expiration:    nullable.NullableTime{},
		Disabled:      false,
	}
}

func _testAddJWTSigningMaterial(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		materialToAdd     *models.JWTSigningMaterial
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:          "GIVEN a valid jwt signing material EXPECT success",
			materialToAdd: &jwtSigningMaterial1,
		},
		{
			name:          "GIVEN another valid jwt signing material EXPECT success",
			materialToAdd: &jwtSigningMaterial2,
		},
		{
			name:          "GIVEN one more jwt wigning material EXPECT success",
			materialToAdd: &jwtSigningMaterial3,
		},
		{
			name:              "GIVEN a jwt signing material with a key id that already exists EXPECT error code jwt signing material key id not unique",
			materialToAdd:     &duplicateJWTSigningMaterial,
			expectedErrorCode: coreerrors.ErrCodeJWTSigningMaterialKeyIDNotUnique,
		},
		{
			name:              "GIVEN a jwt signing material with no key id EXPECT error code jwt signing material key id missing",
			materialToAdd:     &models.JWTSigningMaterial{},
			expectedErrorCode: coreerrors.ErrCodeJWTSigningMaterialKeyIDMissing,
		},
		{
			name: "GIVEN a jwt signing material with no algorithm type EXPECT error code jwt signing material algorithm type missing",
			materialToAdd: &models.JWTSigningMaterial{
				KeyID: uuid.Must(uuid.NewRandom()).String(),
			},
			expectedErrorCode: coreerrors.ErrCodeJWTSigningMaterialAlgorithmTypeMissing,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := jsms.AddJWTSigningMaterial(context.TODO(), logger, tc.materialToAdd, JSMSTestCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.materialToAdd.AuditData.CreatedByID != JSMSTestCreatedBy {
					t.Errorf("\tcreated by id not expected value: got - %s expected - %s", tc.materialToAdd.AuditData.CreatedByID, JSMSTestCreatedBy)
				}
				if tc.materialToAdd.AuditData.CreatedOnDate.IsZero() {
					t.Error("\tcreated date not populated")
				}
			}
		})
	}
}

func _testGetJWTSigningMaterialByKeyID(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		keyID             string
		expectedJSM       models.JWTSigningMaterial
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a key id EXPECT a result",
			keyID:       jwtSigningMaterial1.KeyID,
			expectedJSM: jwtSigningMaterial1,
		},
		{
			name:              "GIVEN a non existant key id EXPECT error code no jwt signing material found",
			keyID:             "not a valid key id",
			expectedErrorCode: coreerrors.ErrCodeNoJWTSigningMaterialFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsm, err := jsms.GetJWTSigningMaterialByKeyID(context.TODO(), logger, tc.keyID, JSMSTestCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if jsm.KeyID != tc.expectedJSM.KeyID {
					t.Errorf("\texpected jsm key id not returned")
				}
			}
		})
	}
}

func _testGetValidJWTSigningMaterialByAlgorithmType(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name                       string
		algorithmType              jwt.JWTSingingAlgorithmFamily
		expectedJwtSigningMaterial []models.JWTSigningMaterial
		expectedErrorCode          string
	}
	testCases := []testCase{
		{
			name:          "GIVEN an algorithm type with two keys in the data store EXPECT the two results",
			algorithmType: jwt.HMAC,
			expectedJwtSigningMaterial: []models.JWTSigningMaterial{
				jwtSigningMaterial1,
				jwtSigningMaterial2,
			},
		},
		{
			name:                       "GIVEN a non existant algorithm type EXPECT no results returned",
			algorithmType:              "not an existing algorithm type",
			expectedJwtSigningMaterial: []models.JWTSigningMaterial{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := jsms.GetValidJWTSigningMaterialByAlgorithmType(context.TODO(), logger, tc.algorithmType, JSMSTestCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numFound := len(result)
				numExpected := len(tc.expectedJwtSigningMaterial)
				if numFound != numExpected {
					t.Errorf("\tnum results not expected: got - %d expected - %d", numFound, numExpected)
				}
				for _, r := range result {
					if r.AlgorithmType != tc.algorithmType {
						t.Errorf("\talgorithm type of returned result not expected: got - %s expcted - %s", r.AlgorithmType, tc.algorithmType)
					}
				}
			}
		})
	}
}
