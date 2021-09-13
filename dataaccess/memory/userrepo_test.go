package memory

import (
	"context"
	"testing"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
)

func TestMemoryUserRepo(t *testing.T) {
	userRepo := NewMemoryUserRepo()
	t.Run("AddUser", func(t *testing.T) {
		testAddUser(t, userRepo)
	})
}

func testAddUser(t *testing.T, userRepo repo.UserRepo) {
	createdBy := "add user test"
	user := models.User{}
	err := userRepo.AddUser(context.TODO(), &user, createdBy)
	if err != nil {
		t.Error(err)
	}
	retreivedUser, ok := users[user.ID]
	if !ok {
		t.Errorf("expected to find user with id: %s", user.ID)
	}
	if retreivedUser.AuditData.CreatedByID != createdBy {
		t.Errorf("expected retreived user to have created by id: '%s' - got '%s'", createdBy, retreivedUser.AuditData.CreatedByID)
	}
	if retreivedUser.AuditData.CreatedOnDate.IsZero() {
		t.Error("expected retreived user to have created on date populated")
	}
}

func testUpdateUser(t *testing.T, userRepo repo.UserRepo) {

}

func testGetUserById(t *testing.T, userRepo repo.UserRepo) {

}

func testGetUserByPrimaryContact(t *testing.T, userRepo repo.UserRepo) {

}

func testGetUserAndContactByPrimaryContact(t *testing.T, userRepo repo.UserRepo) {

}
