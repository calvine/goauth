package service

import (
	"context"
	"testing"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
)

var (
	userServiceTest_ConfirmedUser models.User

	userServiceTest_ConfirmedUserConfirmedContact   models.Contact
	userServiceTest_ConfirmedUserUnconfirmedContact models.Contact

	userServiceTest_UnconfirmedUser models.User

	userServiceTest_UnconfirmedUserConfirmedContact   models.Contact
	userServiceTest_UnconfirmedUserUnconfirmedContact models.Contact

	userServiceTest_UserToRegister models.User

	userServiceTest_UserToRegisterContact models.Contact
)

const (
	userServiceTest_CreatedBy = "user service tests"

	userServiceTest_ConfirmedPrimaryEmail   = "userserviceconprim@email.com"
	userServiceTest_ConfirmedSecondaryEmail = "userserviceconsec@email.com"

	userServiceTest_UnconfirmedPrimaryEmail   = "userserviceunconprim@email.com"
	userServiceTest_UnconfirmedSecondaryEmail = "userserviceunconsec@email.com"

	userServiceTest_UserToRegisterPrimaryEmail = "userserviceunconprim@email.com"
)

func TestUserService(t *testing.T) {
	userService := buildUserService(t)

	t.Run("GetName", func(t *testing.T) {
		_testUserServiceGetName(t, userService)
	})

	t.Run("GetUserByConfirmedContact", func(t *testing.T) {
		_testGetUserByConfirmedContact(t, userService)
	})

	// t.Run("AddUser", func(t *testing.T) {
	// 	_testAddUser(t, userService)
	// })

	// t.Run("UpdateUser", func(t *testing.T) {
	// 	_testUpdateUser(t, userService)
	// })

	t.Run("RegisterUserAndPrimaryContact", func(t *testing.T) {
		_testRegisterUserAndPrimaryContact(t, userService)
	})

	t.Run("GetUserPrimaryContact", func(t *testing.T) {
		_testGetUserPrimaryContact(t, userService)
	})

	t.Run("GetUsersContacts", func(t *testing.T) {
		_testGetUsersContacts(t, userService)
	})

	t.Run("GetUsersConfirmedContacts", func(t *testing.T) {
		_testGetUsersConfirmedContacts(t, userService)
	})

	t.Run("AddContact", func(t *testing.T) {
		_testAddContact(t, userService)
	})

	t.Run("UpdateContact", func(t *testing.T) {
		_testUpdateContact(t, userService)
	})

	t.Run("ConfirmContact", func(t *testing.T) {
		_testConfirmContact(t, userService)
	})
}

func setupTestUserServiceData(t *testing.T, userRepo repo.UserRepo, contactRepo repo.ContactRepo) {
	userServiceTest_ConfirmedUser = models.User{
		PasswordHash: "does not matter",
	}
	err := userRepo.AddUser(context.TODO(), &userServiceTest_ConfirmedUser, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	userServiceTest_ConfirmedUserConfirmedContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, true)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUserConfirmedContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create primary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}
}

func buildUserService(t *testing.T) services.UserService {
	userRepo := memory.NewMemoryUserRepo()
	contactRepo := memory.NewMemoryContactRepo()
	tokenRepo := memory.NewMemoryTokenRepo()
	tokenService := NewTokenService(tokenRepo)
	emailService, _ := NewEmailService(NoOpEmailService, nil)
	userService := NewUserService(userRepo, contactRepo, tokenService, emailService)
	setupTestUserServiceData(t, userRepo, contactRepo)
	return userService
}

func _testUserServiceGetName(t *testing.T, userService services.UserService) {
	serviceName := userService.GetName()
	expectedServiceName := "userService"
	if serviceName != expectedServiceName {
		t.Errorf("service name is not what was expected: got %s - expected %s", serviceName, expectedServiceName)
	}
}

func _testGetUserByConfirmedContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testRegisterUserAndPrimaryContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

// func _testAddUser(t *testing.T, userService services.UserService) {
// t.Error(coreerrors.NewNotImplementedError(true))
// t.Fail()
// }

// func _testUpdateUser(t *testing.T, userService services.UserService) {
// t.Error(coreerrors.NewNotImplementedError(true))
// t.Fail()
// }

func _testGetUserPrimaryContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testGetUsersContacts(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testGetUsersConfirmedContacts(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testAddContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testUpdateContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testConfirmContact(t *testing.T, userService services.UserService) {
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}
