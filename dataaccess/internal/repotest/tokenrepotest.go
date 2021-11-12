package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	richerrors "github.com/calvine/richerror/errors"
)

var (
	testCSRFToken          models.Token
	testPasswordResetToken models.Token
	testSessionToken       models.Token
	testExpiredToken       models.Token
)

const (
	ARBITRARY_DATA_KEY   = "arbitrary_data"
	ARBITRARY_DATA_VALUE = "1234"
)

// TODO: encapsulate sub tests so the run more coherently...

func testTokenRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	_makeTokens(t)
	t.Run("PutToken", func(t *testing.T) {
		_testPutToken(t, *testHarness.TokenRepo)
	})
	t.Run("DeleteToken", func(t *testing.T) {
		_testDeleteToken(t, *testHarness.TokenRepo)
	})
	t.Run("GetToken", func(t *testing.T) {
		_testGetToken(t, *testHarness.TokenRepo)
	})
}

func _makeTokens(t *testing.T) {
	var err richerrors.RichError
	testCSRFToken, err = models.NewToken("", models.TokenTypeCSRF, time.Second*20)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add token of type %s: %s", testCSRFToken.TokenType.String(), err.GetErrorCode())
	}
	testPasswordResetToken, err = models.NewToken("", models.TokenTypePasswordReset, time.Second*20)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add token of type %s: %s", testPasswordResetToken.TokenType.String(), err.GetErrorCode())
	}
	testSessionToken, err = models.NewToken("fake_user_id", models.TokenTypeSession, time.Second*20)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add token of type %s: %s", testSessionToken.TokenType.String(), err.GetErrorCode())
	}
	testSessionToken.AddMetaData(ARBITRARY_DATA_KEY, ARBITRARY_DATA_VALUE)
	testExpiredToken, err = models.NewToken("fake_user_id2", models.TokenTypeSession, time.Second*-20)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add token of type %s: %s", testExpiredToken.TokenType.String(), err.GetErrorCode())
	}
}

func _testPutToken(t *testing.T, tokenRepo repo.TokenRepo) {
	err := tokenRepo.PutToken(context.TODO(), testCSRFToken)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to put token of type %s: %s", testCSRFToken.TokenType.String(), err.GetErrorCode())
	}
	err = tokenRepo.PutToken(context.TODO(), testPasswordResetToken)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to put token of type %s: %s", testPasswordResetToken.TokenType.String(), err.GetErrorCode())
	}
	err = tokenRepo.PutToken(context.TODO(), testSessionToken)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to put token of type %s: %s", testSessionToken.TokenType.String(), err.GetErrorCode())
	}
	err = tokenRepo.PutToken(context.TODO(), testExpiredToken)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to put token of type %s: %s", testExpiredToken.TokenType.String(), err.GetErrorCode())
	}
}

func _testDeleteToken(t *testing.T, tokenRepo repo.TokenRepo) {
	err := tokenRepo.DeleteToken(context.TODO(), testPasswordResetToken.Value)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to delete token with type %s: %s", testPasswordResetToken.TokenType.String(), err.GetErrorCode())
	}
}

func _testGetToken(t *testing.T, tokenRepo repo.TokenRepo) {
	_, err := tokenRepo.GetToken(context.TODO(), testPasswordResetToken.Value)
	if err.GetErrorCode() != errors.ErrCodeTokenNotFound {
		t.Error("testPasswordResetToken found inspite of being deleted in the tokenRepo.DeleteToken test...")
	}
	expiredToken, err := tokenRepo.GetToken(context.TODO(), testExpiredToken.Value)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get testExpiredToken from repo: %s", err.GetErrorCode())
	}
	if !expiredToken.IsExpired() {
		t.Errorf("testExpiredToken should be expired but expiration is: %s and it is now %s", expiredToken.Expiration.UTC().String(), time.Now().UTC().String())
	}
	csrfToken, err := tokenRepo.GetToken(context.TODO(), testCSRFToken.Value)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get testCSRFToken from repo: %s", err.GetErrorCode())
	}
	if csrfToken.TokenType != models.TokenTypeCSRF {
		t.Errorf("testCSRFToken expected token type does not match expected value: got: %s - expected: %s", csrfToken.TokenType.String(), testCSRFToken.TokenType.String())
	}
	if csrfToken.IsExpired() {
		t.Errorf("testCSRFToken should not be expired but expiration is: %s and it is now %s", expiredToken.Expiration.UTC().String(), time.Now().UTC().String())
	}
	sessionToken, err := tokenRepo.GetToken(context.TODO(), testSessionToken.Value)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to get testSessionToken from repo: %s", err.GetErrorCode())
	}
	if sessionToken.TokenType != models.TokenTypeSession {
		t.Errorf("testSessionToken expected token type does not match expected value: got: %s - expected: %s", sessionToken.TokenType.String(), testSessionToken.TokenType.String())
	}
	if sessionToken.IsExpired() {
		t.Errorf("testSessionToken should not be expired but expiration is: %s and it is now %s", expiredToken.Expiration.UTC().String(), time.Now().UTC().String())
	}
	if sessionToken.TargetID != testSessionToken.TargetID {
		t.Errorf("testSessionToken expected target id does not match expected value: got: %s - expected: %s", sessionToken.TargetID, testSessionToken.TargetID)
	}
	if len(sessionToken.MetaData) == 0 {
		t.Error("testSessionToken metadata is not present")
	}
	arbitraryValue, ok := sessionToken.MetaData[ARBITRARY_DATA_KEY]
	if !ok {
		t.Errorf("testSessionToken metdata map does not contain key %s", ARBITRARY_DATA_KEY)
	}
	if arbitraryValue != ARBITRARY_DATA_VALUE {
		t.Errorf("testSessionToken expected metadata for key %s does not match expected value: got: %s - expected: %s", ARBITRARY_DATA_KEY, arbitraryValue, ARBITRARY_DATA_VALUE)
	}
}
