package service

import (
	"testing"

	"github.com/calvine/goauth/core/models"
	coreservices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
)

var (
	jwtSigningMaterial1 models.JWTSigningMaterial
	jwtSigningMaterial2 models.JWTSigningMaterial
	jwtSigningMaterial3 models.JWTSigningMaterial
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
	t.Error("TODO add the test data and implement tests")
}

func _testAddJWTSigningMaterial(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	type testCase struct {
		name              string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {})
	}
	t.Error("\ttest not implemented")
}

func _testGetJWTSigningMaterialByKeyID(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	type testCase struct {
		name              string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {})
	}
	t.Error("\ttest not implemented")
}

func _testGetValidJWTSigningMaterialByAlgorithmType(t *testing.T, jsms coreservices.JWTSigningMaterialService) {
	type testCase struct {
		name              string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {})
	}
	t.Error("\ttest not implemented")
}
