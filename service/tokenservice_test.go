package service

import (
	"context"
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
	"github.com/calvine/richerror/errors"
)

var (
	csrfToken                 models.Token
	confirmContactToken       models.Token
	passwordResetToken        models.Token
	sessionToken              models.Token
	shortDurationSessionToken models.Token
	expiredSessionToken       models.Token
)

const (
	tokenUserID string = "test_token_user_id"

	testKey   = "testkey"
	testValue = "testvalue"

	sesssionTokenExpirationDuration time.Duration = time.Millisecond * 200
)

func TestTokenService(t *testing.T) {
	tokenRepo := memory.NewMemoryTokenRepo()
	tokenService := NewTokenService(tokenRepo)

	setupTestTokens(t)

	// put token
	t.Run("PutToken", func(t *testing.T) {
		_testPutToken(t, tokenService)
	})

	// get token
	t.Run("GetToken", func(t *testing.T) {
		_testGetToken(t, tokenService)
	})

	// delete token
	t.Run("DeleteToken", func(t *testing.T) {
		_testDeleteToken(t, tokenService)
	})
}

func setupTestTokens(t *testing.T) {
	var err errors.RichError
	csrfToken, err = models.NewToken(tokenUserID, models.TokenTypeCSRF, time.Minute*1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	confirmContactToken, err = models.NewToken(tokenUserID, models.TokenTypeConfirmContact, time.Minute*1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	passwordResetToken, err = models.NewToken(tokenUserID, models.TokenTypePasswordReset, time.Minute*1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	sessionToken, err = models.NewToken(tokenUserID, models.TokenTypeSession, time.Minute*1)
	sessionToken.AddMetaData(testKey, testValue)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	shortDurationSessionToken, err = models.NewToken(tokenUserID, models.TokenTypeSession, sesssionTokenExpirationDuration)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	expiredSessionToken, err = models.NewToken(tokenUserID, models.TokenTypeSession, time.Minute*-1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
}

func _testPutToken(t *testing.T, tokenService services.TokenService) {
	// success
	t.Run("Success", func(t *testing.T) {
		__testPutTokenSuccess(t, tokenService)
	})

	// failure empty token value
	t.Run("Failure empty token value", func(t *testing.T) {
		__testPutTokenFailureEmptyTokenValue(t, tokenService)
	})

	// failure expiration in past
	t.Run("Failure expiration in past", func(t *testing.T) {
		__testPutTokenFailureExpiredToken(t, tokenService)
	})

	// failure token type invalid
	t.Run("Failure token type in past", func(t *testing.T) {
		__testPutTokenFailureInvalidTokenType(t, tokenService)
	})
}

func __testPutTokenSuccess(t *testing.T, tokenService services.TokenService) {
	err := tokenService.PutToken(context.TODO(), csrfToken)
	if err != nil {
		t.Errorf("failed to add token %v got error %s: %s", csrfToken, err.GetErrorCode(), err.Error())
	}
	err = tokenService.PutToken(context.TODO(), confirmContactToken)
	if err != nil {
		t.Errorf("failed to add token %v got error %s: %s", csrfToken, err.GetErrorCode(), err.Error())
	}
	err = tokenService.PutToken(context.TODO(), passwordResetToken)
	if err != nil {
		t.Errorf("failed to add token %v got error %s: %s", csrfToken, err.GetErrorCode(), err.Error())
	}
	err = tokenService.PutToken(context.TODO(), sessionToken)
	if err != nil {
		t.Errorf("failed to add token %v got error %s: %s", csrfToken, err.GetErrorCode(), err.Error())
	}
	err = tokenService.PutToken(context.TODO(), shortDurationSessionToken)
	if err != nil {
		t.Errorf("failed to add token %v got error %s: %s", csrfToken, err.GetErrorCode(), err.Error())
	}
}

func __testPutTokenFailureEmptyTokenValue(t *testing.T, tokenService services.TokenService) {
	emptyValueToken, err := models.NewToken("", models.TokenTypeCSRF, time.Minute*1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	emptyValueToken.Value = ""
	err = tokenService.PutToken(context.TODO(), emptyValueToken)
	if err == nil {
		t.Error("expected error because token expiration has passed")
	}
	if err.GetErrorCode() != coreerrors.ErrCodeMalfomedToken {
		t.Errorf("expected malformed token error bug got %s: %s", err.GetErrorCode(), err.Error())
	}
}

func __testPutTokenFailureExpiredToken(t *testing.T, tokenService services.TokenService) {
	err := tokenService.PutToken(context.TODO(), expiredSessionToken)
	if err == nil {
		t.Error("expected error because token expiration has passed")
	}
	if err.GetErrorCode() != coreerrors.ErrCodeMalfomedToken {
		t.Errorf("expected malformed token error bug got %s: %s", err.GetErrorCode(), err.Error())
	}
}

func __testPutTokenFailureInvalidTokenType(t *testing.T, tokenService services.TokenService) {
	invalidToken, err := models.NewToken("", models.TokenTypeInvalid, time.Minute*1)
	if err != nil {
		t.Errorf("faiiled to create token for testing token service: %s", err.Error())
	}
	err = tokenService.PutToken(context.TODO(), invalidToken)
	if err == nil {
		t.Error("expected error because token expiration has passed")
	}
	if err.GetErrorCode() != coreerrors.ErrCodeMalfomedToken {
		t.Errorf("expected malformed token error bug got %s: %s", err.GetErrorCode(), err.Error())
	}
}

func _testGetToken(t *testing.T, tokenService services.TokenService) {
	// success
	t.Run("Success", func(t *testing.T) {
		__testGetTokenSuccess(t, tokenService)
	})
	// success with meta data
	t.Run("Success with metadata", func(t *testing.T) {
		__testGetTokenSuccessWithMetadata(t, tokenService)
	})
	// failure token expired
	t.Run("Failure token expired", func(t *testing.T) {
		__testGetTokenFailureTokenExpired(t, tokenService)
	})
	// failure wrong token type provided
	t.Run("Failure wrong token type provided", func(t *testing.T) {
		__testGetTokenFailureWrongTokenType(t, tokenService)
	})
}

func __testGetTokenSuccess(t *testing.T, tokenService services.TokenService) {
	token, err := tokenService.GetToken(context.TODO(), csrfToken.Value, models.TokenTypeCSRF)
	if err != nil {
		t.Errorf("failed to get token got error %s: %s", err.GetErrorCode(), err.Error())
	}
	if token.Value != csrfToken.Value {
		t.Errorf("retreived token value does not match expected value: got: %s - expected %s", token.Value, csrfToken.Value)
	}
	if token.TargetID != csrfToken.TargetID {
		t.Errorf("retreived token target id does not match expected value: got: %s - expected %s", token.TargetID, csrfToken.TargetID)
	}
}

func __testGetTokenSuccessWithMetadata(t *testing.T, tokenService services.TokenService) {
	token, err := tokenService.GetToken(context.TODO(), sessionToken.Value, models.TokenTypeSession)
	if err != nil {
		t.Errorf("failed to get token got error %s: %s", err.GetErrorCode(), err.Error())
	}
	if token.Value != sessionToken.Value {
		t.Errorf("retreived token value does not match expected value: got: %s - expected %s", token.Value, sessionToken.Value)
	}
	if token.TargetID != sessionToken.TargetID {
		t.Errorf("retreived token target id does not match expected value: got: %s - expected %s", token.TargetID, sessionToken.TargetID)
	}
	if token.MetaData == nil {
		t.Error("")
	}
	if token.MetaData[testKey] != testValue {
		t.Errorf("retreived token metadata for key %s was not the expected value: got %s - expected %s", testKey, token.MetaData[testKey], testValue)
	}
}

func __testGetTokenFailureTokenExpired(t *testing.T, tokenService services.TokenService) {
	time.Sleep(sesssionTokenExpirationDuration)
	_, err := tokenService.GetToken(context.TODO(), shortDurationSessionToken.Value, models.TokenTypeSession)
	if err == nil {
		t.Error("expected error because the token is expired")
	}
	if err.GetErrorCode() != coreerrors.ErrCodeExpiredToken {
		t.Errorf("expected expired token error but got %s: %s", err.GetErrorCode(), err.Error())
	}
}

func __testGetTokenFailureWrongTokenType(t *testing.T, tokenService services.TokenService) {
	_, err := tokenService.GetToken(context.TODO(), sessionToken.Value, models.TokenTypeCSRF)
	if err == nil {
		t.Error("expected error because the token is expired")
	}
	if err.GetErrorCode() != coreerrors.ErrCodeWrongTokenType {
		t.Errorf("expected expired token error but got %s: %s", err.GetErrorCode(), err.Error())
	}
}

func _testDeleteToken(t *testing.T, tokenService services.TokenService) {
	// success
	t.Run("Success", func(t *testing.T) {
		__testDeleteTokenSuccess(t, tokenService)
	})
	//Do we need this test?
	// // failure token not found
	// t.Run("Failure token not found", func(t *testing.T) {
	// 	__testDeleteTokenFailuireTokenNotFound(t, tokenService)
	// })
}

func __testDeleteTokenSuccess(t *testing.T, tokenService services.TokenService) {
	err := tokenService.DeleteToken(context.TODO(), csrfToken.Value)
	if err != nil {
		t.Errorf("failed to delete token got error %s: %s", err.GetErrorCode(), err.Error())
	}
}

// func __testDeleteTokenFailuireTokenNotFound(t *testing.T, tokenService services.TokenService) {
// 	err := tokenService.DeleteToken(context.TODO(), "not a real token value9867568797567687")
// 	if err != nil {
// 		t.Errorf("failed to delete token got error %s: %s", err.GetErrorCode(), err.Error())
// 	}
// }
