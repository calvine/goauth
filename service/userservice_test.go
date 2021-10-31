package service

import (
	"testing"

	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
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

func buildUserService(t *testing.T) services.UserService {
	userRepo := memory.NewMemoryUserRepo()
	contactRepo := memory.NewMemoryContactRepo()
	tokenRepo := memory.NewMemoryTokenRepo()
	tokenService := NewTokenService(tokenRepo)
	emailService, _ := NewEmailService(NoOpEmailService, nil)
	userService := NewUserService(userRepo, contactRepo, tokenService, emailService)
	// TODO: set up test data.
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

}

// func _testAddUser(t *testing.T, userService services.UserService) {

// }

// func _testUpdateUser(t *testing.T, userService services.UserService) {

// }

func _testGetUserPrimaryContact(t *testing.T, userService services.UserService) {

}

func _testGetUsersContacts(t *testing.T, userService services.UserService) {

}

func _testGetUsersConfirmedContacts(t *testing.T, userService services.UserService) {

}

func _testAddContact(t *testing.T, userService services.UserService) {

}

func _testUpdateContact(t *testing.T, userService services.UserService) {

}

func _testConfirmContact(t *testing.T, userService services.UserService) {

}
