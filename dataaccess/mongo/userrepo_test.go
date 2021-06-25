package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core/models"
)

var (
	testUser1 = models.User{
		Id:           "uniqueId1",
		PasswordHash: "passwordhash1",
		Salt:         "salt1",
	}
	testUser2 = models.User{
		Id:           "uniqueId2",
		PasswordHash: "passwordhash2",
		Salt:         "salt2",
	}
)

func testMongoUserRepo(t *testing.T, userRepo *userRepo) {

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
		t.Fail()
	})

}

func _testAddUser(t *testing.T, userRepo *userRepo) {
	createdById := "test1"

	err := userRepo.AddUser(context.TODO(), &testUser1, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser1.AuditData.CreatedById != createdById {
		t.Error("failed to set the users CreatedByID to the right value", testUser1.AuditData.CreatedById, createdById)
	}

	createdById = "test2"

	err = userRepo.AddUser(context.TODO(), &testUser2, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser2.AuditData.CreatedById != createdById {
		t.Error("failed to set the users CreatedByID to the right value", testUser2.AuditData.CreatedById, createdById)
	}
}

func _testUpdateUser(t *testing.T, userRepo *userRepo) {
	preUpdateDate := time.Now().UTC()
	newPasswordHash := "another secure password hash"
	newSalt := "change password = change salt"
	testUser1.PasswordHash = newPasswordHash
	testUser1.Salt = newSalt
	err := userRepo.UpdateUser(context.TODO(), &testUser1, testUser1.Id)
	if err != nil {
		t.Error("failed to update user", err)
	}
	if !testUser1.AuditData.ModifiedOnDate.Value.After(preUpdateDate) {
		t.Error("ModifiedOnDate should be after the preUpdate for test", preUpdateDate, testUser1.AuditData.ModifiedOnDate)
	}
	retreivedUser, err := userRepo.GetUserById(context.TODO(), testUser1.Id)
	if err != nil {
		t.Error("failed to retreive updated user to check that fields were updated.", err)
	}
	if retreivedUser.PasswordHash != newPasswordHash {
		t.Error("password hash did not update.", retreivedUser.PasswordHash, newPasswordHash)
	}
	if retreivedUser.Salt != newSalt {
		t.Error("salt did not update.", retreivedUser.Salt, newSalt)
	}
}

func _testGetUserById(t *testing.T, userRepo *userRepo) {
	userId := testUser2.Id
	retreivedUser, err := userRepo.GetUserById(context.TODO(), userId)
	if err != nil {
		t.Error("error getting user with id", userId, err)
	}
	if retreivedUser.PasswordHash != testUser2.PasswordHash || retreivedUser.Salt != testUser2.Salt {
		t.Error("retreivedUser should have same data as user with id tested", retreivedUser, testUser2)
	}
}
