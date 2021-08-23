package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
)

var (
	testUser1 = models.User{
		PasswordHash: "passwordhash1",
	}
	// TODO: Finish implementing for get user my password reset token test.
	testUser2 = models.User{
		PasswordResetToken: nullable.NullableString{HasValue: true, Value: "passwordresettoken1"},
		// non expired token
		PasswordResetTokenExpiration: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * 15).UTC()},
	}
	testUser3 = models.User{
		PasswordResetToken: nullable.NullableString{HasValue: true, Value: "passwordresettoken2"},
		// expired token
		PasswordResetTokenExpiration: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -15)},
	}
)

func testMongoUserRepo(t *testing.T, userRepo userRepo) {
	// functionality tests
	t.Run("userRepo.AddUser", func(t *testing.T) {
		_testAddUser(t, userRepo)
	})
	t.Run("userRepo.UpdateUser", func(t *testing.T) {
		_testUpdateUser(t, userRepo)
	})
	t.Run("userRepo.GetUserById", func(t *testing.T) {
		_testGetUserById(t, userRepo)
	})
	t.Run("userRepo.GetUserByPrimaryContact", func(t *testing.T) {
		_testGetUserByPrimaryContact(t, userRepo)
	})
	t.Run("userRepo.GetUserAndContactByPrimaryContact", func(t *testing.T) {
		_testGetUserAndContactByPrimaryContact(t, userRepo)
	})
	t.Run("userRepo.GetUserByPasswordResetToken", func(t *testing.T) {
		_testGetUserByPasswordResetToken(t, userRepo)
	})
}

func _testAddUser(t *testing.T, userRepo userRepo) {
	createdById := "user repos tests"

	err := userRepo.AddUser(context.TODO(), &testUser1, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser1.AuditData.CreatedByID != createdById {
		t.Error("failed to set the test user 1 CreatedByID to the right value", testUser1.AuditData.CreatedByID, createdById)
	}

	err = userRepo.AddUser(context.TODO(), &testUser2, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser2.AuditData.CreatedByID != createdById {
		t.Error("failed to set the test user 2 CreatedByID to the right value", testUser2.AuditData.CreatedByID, createdById)
	}

	err = userRepo.AddUser(context.TODO(), &testUser3, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser3.AuditData.CreatedByID != createdById {
		t.Error("failed to set the test user 3 CreatedByID to the right value", testUser3.AuditData.CreatedByID, createdById)
	}
}

func _testUpdateUser(t *testing.T, userRepo userRepo) {
	preUpdateDate := time.Now().UTC()
	newPasswordHash := "another secure password hash"
	testUser1.PasswordHash = newPasswordHash
	err := userRepo.UpdateUser(context.TODO(), &testUser1, testUser1.ID)
	if err != nil {
		t.Error("failed to update user", err)
	}
	if !testUser1.AuditData.ModifiedOnDate.Value.After(preUpdateDate) {
		t.Error("ModifiedOnDate should be after the preUpdate for test", preUpdateDate, testUser1.AuditData.ModifiedOnDate)
	}
	retreivedUser, err := userRepo.GetUserById(context.TODO(), testUser1.ID)
	if err != nil {
		t.Error("failed to retreive updated user to check that fields were updated.", err)
	}
	if retreivedUser.PasswordHash != newPasswordHash {
		t.Error("password hash did not update.", retreivedUser.PasswordHash, newPasswordHash)
	}
}

func _testGetUserById(t *testing.T, userRepo userRepo) {
	userId := initialTestUser.ID
	retreivedUser, err := userRepo.GetUserById(context.TODO(), userId)
	if err != nil {
		t.Error("error getting user with id", userId, err)
	}
	if retreivedUser.PasswordHash != initialTestUser.PasswordHash {
		t.Error("retreivedUser should have same data as user with id tested", retreivedUser, initialTestUser)
	}
}

func _testGetUserByPrimaryContact(t *testing.T, userRepo userRepo) {
	contactType, principal := core.CONTACT_TYPE_EMAIL, "InitialTestUser@email.com"
	retreivedUser, err := userRepo.GetUserByPrimaryContact(context.TODO(), contactType, principal)
	if err != nil {
		t.Error("failed to retreive user via primary contact info", contactType, principal, err)
	}
	if retreivedUser.ID != initialTestUser.ID {
		t.Error("expected retreivedUser and initialTestUser Id to match", retreivedUser.ID, initialTestUser.ID)
	}
}

func _testGetUserAndContactByPrimaryContact(t *testing.T, userRepo userRepo) {
	contactType, principal := core.CONTACT_TYPE_EMAIL, "InitialTestUser@email.com"
	retreivedUser, retreivedContact, err := userRepo.GetUserAndContactByPrimaryContact(context.TODO(), contactType, principal)
	if err != nil {
		t.Error("failed to retreive user via primary contact info", contactType, principal, err)
	}
	if retreivedUser.ID != initialTestUser.ID {
		t.Error("expected retreivedUser and initialTestUser Id to match", retreivedUser.ID, initialTestUser.ID)
	}
	if retreivedContact.Principal != principal {
		t.Error("expected retreivedContact.Principal and the test principal to be the same", retreivedContact.Principal, principal)
	}
	if retreivedContact.UserID != retreivedUser.ID {
		t.Error("expected user.ID and contact.userID to match", retreivedContact.UserID, retreivedUser.ID)
	}
}

func _testGetUserByPasswordResetToken(t *testing.T, userRepo userRepo) {
	retreivedUser, err := userRepo.GetUserByPasswordResetToken(context.TODO(), testUser2.PasswordResetToken.Value)
	if err != nil {
		t.Error("failed to retreive user via password reset token", testUser2.PasswordResetToken.Value, err)
	}
	if retreivedUser.ID != testUser2.ID {
		t.Error("expected retreivedUser and testUser2 Id to match", retreivedUser.ID, testUser2.ID)
	}
	if retreivedUser.PasswordResetToken.Value != testUser2.PasswordResetToken.Value {
		t.Error("expected retreivedUser.PasswordResetToken and testUser2.PasswordResetToken to be the same", retreivedUser.PasswordResetToken.Value, testUser2.PasswordResetToken.Value)
	}
	// NOTE: commented out for now because mongodb does not record micro seconds...
	// if retreivedUser.PasswordResetTokenExpiration.Value != testUser2.PasswordResetTokenExpiration.Value {
	// 	t.Error("expected retreivedUser.PasswordResetTokenExpiration and testUser2.PasswordResetTokenExpiration to match", retreivedUser.PasswordResetTokenExpiration.Value, testUser2.PasswordResetTokenExpiration.Value)
	// }
}
