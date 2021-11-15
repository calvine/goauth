package service

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/dataaccess/memory"
	"github.com/calvine/goauth/internal/testutils"
	"go.uber.org/zap/zaptest"
)

var (
	userServiceTest_ConfirmedUser models.User

	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact     models.Contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact   models.Contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact models.Contact

	userServiceTest_UnconfirmedUser models.User

	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact models.Contact

	userServiceTest_UserToRegister models.User

	userServiceTest_UserToRegisterContact models.Contact
)

const (
	userServiceTest_CreatedBy = "user service tests"

	userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail     = "userserviceconprim@email.com"
	userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail   = "userserviceconsec@email.com"
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail = "userserviceunconsec@email.com"

	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryEmail = "userserviceunconprim@email.com"

	userServiceTest_UserToRegisterEmail = "userservicetoregister@email.com"
)

func TestUserService(t *testing.T) {
	userService := buildUserService(t)

	t.Run("GetName", func(t *testing.T) {
		_testUserServiceGetName(t, userService)
	})

	t.Run("GetUserAndContactByConfirmedContact", func(t *testing.T) {
		_testGetUserAndContactByConfirmedContact(t, userService)
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

	// t.Run("GetUserPrimaryContact", func(t *testing.T) {
	// 	_testGetUserPrimaryContact(t, userService)
	// })

	// t.Run("GetUsersContacts", func(t *testing.T) {
	// 	_testGetUsersContacts(t, userService)
	// })

	// t.Run("GetUsersConfirmedContacts", func(t *testing.T) {
	// 	_testGetUsersConfirmedContacts(t, userService)
	// })

	// t.Run("AddContact", func(t *testing.T) {
	// 	_testAddContact(t, userService)
	// })

	// t.Run("UpdateContact", func(t *testing.T) {
	// 	_testUpdateContact(t, userService)
	// })

	// t.Run("ConfirmContact", func(t *testing.T) {
	// 	_testConfirmContact(t, userService)
	// })
}

func setupTestUserServiceData(t *testing.T, userRepo repo.UserRepo, contactRepo repo.ContactRepo) {
	// add user with confirmed contact
	userServiceTest_ConfirmedUser = models.User{
		PasswordHash: "does not matter",
	}
	err := userRepo.AddUser(context.TODO(), &userServiceTest_ConfirmedUser, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed primary contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, true)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedPrimaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create confirmedprimary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed seconday contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, false)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user unconfirmed seconday contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, false)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add unconfirmed user
	userServiceTest_UnconfirmedUser = models.User{
		PasswordHash: "does not matter",
	}
	err = userRepo.AddUser(context.TODO(), &userServiceTest_UnconfirmedUser, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create user with unconfirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add unconfrimed user unconfirmed primary contact
	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact = models.NewContact(userServiceTest_UnconfirmedUser.ID, "", userServiceTest_UnconfirmedUser_UnconfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, true)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create unconfirmed primary contact for user with no confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}
}

func buildUserService(t *testing.T) services.UserService {
	users := make(map[string]models.User)
	contacts := make(map[string]models.Contact)
	userRepo, err := memory.NewMemoryUserRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	contactRepo, err := memory.NewMemoryContactRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
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

func _testGetUserAndContactByConfirmedContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		contactPrincipal  string
		contactType       string
		expectedUserID    string
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "Given confirmed primary contact Return User",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedUserID:    userServiceTest_ConfirmedUser.ID,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
		},
		{
			name:              "Given confirmed secondary contact Return User",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedUserID:    userServiceTest_ConfirmedUser.ID,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
		},
		{
			name:              "Given unconfirmed contact Return error code no user found",
			contactPrincipal:  userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
		{
			name:              "Given non existant contact Return No User Found Error",
			contactPrincipal:  "ojhgfiujwsfiogh@oiwujhgfiwsrb.dofuhsfoiuds",
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, contact, err := userService.GetUserAndContactByConfirmedContact(context.TODO(), logger, tc.contactType, tc.contactPrincipal, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if user.ID != tc.expectedUserID {
					t.Errorf("\tuser id did not match expected: got - %s expected - %s", user.ID, tc.expectedUserID)
					t.Fail()
				}
				if contact.ID != tc.expectedContactID {
					t.Errorf("\tcontact id did not match expected: got - %s expected - %s", contact.ID, tc.expectedContactID)
					t.Fail()
				}
			}
		})
	}
}

// // func _testAddUser(t *testing.T, userService services.UserService) {
// // t.Error(coreerrors.NewNotImplementedError(true))
// // t.Fail()
// // }

// // func _testUpdateUser(t *testing.T, userService services.UserService) {
// // t.Error(coreerrors.NewNotImplementedError(true))
// // t.Fail()
// // }

func _testRegisterUserAndPrimaryContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		contactPrincipal  string
		contactType       string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:             "Given unregistered contact Return Successful registration new user and contact",
			contactPrincipal: userServiceTest_UserToRegisterEmail,
			contactType:      core.CONTACT_TYPE_EMAIL,
		},
		{
			name:             "Given previously registered unconfirmed contact Return successful registration new user and contact",
			contactPrincipal: userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail,
			contactType:      core.CONTACT_TYPE_EMAIL,
		},
		{
			name:              "Given the provided contact is already confirmed in the data store Return error contact already confirmed",
			contactPrincipal:  userServiceTest_UserToRegisterEmail,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeRegistrationContactAlreadyRegistered,
		},
		// TODO: create test case for multiple confirmed instances of a contact returning the appropriate error...
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.RegisterUserAndPrimaryContact(context.TODO(), logger, tc.contactType, tc.contactPrincipal, userServiceTest_CreatedBy)
			if err != nil {
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
			} else {

			}
		})
	}
}

// func _testGetUserPrimaryContact(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }

// func _testGetUsersContacts(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }

// func _testGetUsersConfirmedContacts(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }

// func _testAddContact(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }

// func _testUpdateContact(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }

// func _testConfirmContact(t *testing.T, userService services.UserService) {
// 	t.Error(coreerrors.NewNotImplementedError(true))
// 	t.Fail()
// }
