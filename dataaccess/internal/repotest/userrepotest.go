package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

var (
	testUser1 = models.User{
		PasswordHash: "passwordhash1",
	}
	nonExistantUserID string
)

func setupUserRepoTestData(_ *testing.T, testingHarness RepoTestHarnessInput) {
	nonExistantUserID = testingHarness.IDGenerator(false)
}

func testUserRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	setupUserRepoTestData(t, testHarness)
	// functionality tests
	t.Run("AddUser", func(t *testing.T) {
		_testAddUser(t, *testHarness.UserRepo)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		_testUpdateUser(t, *testHarness.UserRepo)
	})
	t.Run("GetUserByID", func(t *testing.T) {
		_testGetUserByID(t, *testHarness.UserRepo)
	})
	t.Run("GetUserByPrimaryContact", func(t *testing.T) {
		_testGetUserByPrimaryContact(t, *testHarness.UserRepo)
	})
	t.Run("GetUserAndContactByConfirmedContact", func(t *testing.T) {
		_testGetUserAndContactByConfrimedContact(t, *testHarness.UserRepo)
	})
}

func _testAddUser(t *testing.T, userRepo repo.UserRepo) {
	createdByID := "user repos tests"

	err := userRepo.AddUser(context.TODO(), &testUser1, createdByID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to add user to database", err.GetErrorCode())
	}
	if testUser1.AuditData.CreatedByID != createdByID {
		t.Error("failed to set the test user 1 CreatedByID to the right value", testUser1.AuditData.CreatedByID, createdByID)
	}
}

func _testUpdateUser(t *testing.T, userRepo repo.UserRepo) {
	preUpdateDate := time.Now().UTC()
	newPasswordHash := "another secure password hash"
	testUser1.PasswordHash = newPasswordHash
	err := userRepo.UpdateUser(context.TODO(), &testUser1, testUser1.ID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to update user", err.GetErrorCode())
	}
	if !testUser1.AuditData.ModifiedOnDate.Value.After(preUpdateDate) {
		t.Error("ModifiedOnDate should be after the preUpdate for test", preUpdateDate, testUser1.AuditData.ModifiedOnDate)
	}
	retreivedUser, err := userRepo.GetUserByID(context.TODO(), testUser1.ID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to retreive updated user to check that fields were updated.", err.GetErrorCode())
	}
	if retreivedUser.PasswordHash != newPasswordHash {
		t.Error("password hash did not update.", retreivedUser.PasswordHash, newPasswordHash)
	}
}

func _testGetUserByID(t *testing.T, userRepo repo.UserRepo) {
	userID := initialTestUser.ID
	retreivedUser, err := userRepo.GetUserByID(context.TODO(), userID)
	if err != nil {
		t.Log(err.Error())
		t.Error("error getting user with id", userID, err.GetErrorCode())
	}
	if retreivedUser.PasswordHash != initialTestUser.PasswordHash {
		t.Error("retreivedUser should have same data as user with id tested", retreivedUser, initialTestUser)
	}
}

func _testGetUserByPrimaryContact(t *testing.T, userRepo repo.UserRepo) {
	contactType, principal := core.CONTACT_TYPE_EMAIL, initialTestConfirmedPrimaryContact.Principal
	retreivedUser, err := userRepo.GetUserByPrimaryContact(context.TODO(), contactType, principal)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to retreive user via primary contact info", contactType, principal, err.GetErrorCode())
	}
	if retreivedUser.ID != initialTestUser.ID {
		t.Error("expected retreivedUser and initialTestUser ID to match", retreivedUser.ID, initialTestUser.ID)
	}
}

func _testGetUserAndContactByConfrimedContact(t *testing.T, userRepo repo.UserRepo) {
	type testCase struct {
		name              string
		contactPrincipal  string
		contactType       string
		expectedContactID string
		expectedUserID    string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN a confirmed contact principal and type EXPECT user and contact to be returned",
			contactPrincipal:  initialTestConfirmedPrimaryContact.Principal,
			contactType:       initialTestConfirmedPrimaryContact.Type,
			expectedContactID: initialTestConfirmedPrimaryContact.ID,
			expectedUserID:    initialTestConfirmedPrimaryContact.UserID,
		},
		{
			name:              "GIVEN a unconfirmed contact principal and type EXPECT error code no user found",
			contactPrincipal:  initialTestUnconfirmedContact.Principal,
			contactType:       initialTestUnconfirmedContact.Type,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
		{
			name:              "GIVEN a non existant contact principal and a valid contact type EXPECT error code no user found",
			contactPrincipal:  "NOT_A_REAL_EMAIL_456yhgtyTUHG@email.org",
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			retreivedUser, retreivedContact, err := userRepo.GetUserAndContactByConfirmedContact(context.TODO(), tc.contactType, tc.contactPrincipal)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if tc.expectedUserID != retreivedUser.ID {
					t.Errorf("\tuser id expected: got - %s expected - %s", retreivedUser.ID, tc.expectedUserID)
					t.Fail()
				}
				if tc.expectedContactID != retreivedContact.ID {
					t.Errorf("\tcontact id expected: got - %s expected - %s", retreivedContact.ID, tc.expectedContactID)
					t.Fail()
				}
			}
		})
	}
}
