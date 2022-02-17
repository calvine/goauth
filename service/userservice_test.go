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
	userServiceTest_EmailService services.EmailService

	userServiceTest_ConfirmedUser models.User

	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact     models.Contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact   models.Contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact models.Contact

	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact     models.Contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact models.Contact

	userServiceTest_UnconfirmedUser models.User

	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact models.Contact

	userServiceText_ContactRepo repo.ContactRepo
	userServiceText_TokenRepo   repo.TokenRepo
)

const (
	userServiceTest_CreatedBy = "user service tests"

	userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail     = "userserviceconprim@email.com"
	userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail   = "userserviceconsec@email.com"
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail = "userserviceunconsec@email.com"

	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobile     = "111-111-1111"
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobile = "333-333-3333"

	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryEmail = "userserviceunconprim@email.com"

	userServiceTest_UserToRegisterEmail = "userservicetoregister@email.com"
)

func TestUserService(t *testing.T) {
	userService := buildUserService(t)

	t.Run("GetName", func(t *testing.T) {
		_testUserServiceGetName(t, userService)
	})

	t.Run("GetUser", func(t *testing.T) {
		_testUserServiceGetUser(t, userService)
	})

	t.Run("GetUserAndContactByConfirmedContact", func(t *testing.T) {
		_testGetUserAndContactByConfirmedContact(t, userService)
	})

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

	t.Run("GetUsersContactsOfType", func(t *testing.T) {
		_testGetUsersContactsOfType(t, userService)
	})

	t.Run("GetUsersConfirmedContactsOfType", func(t *testing.T) {
		_testGetUsersConfirmedContactsOfType(t, userService)
	})

	t.Run("AddContact", func(t *testing.T) {
		_testAddContact(t, userService)
	})

	t.Run("SetContactAsPrimary", func(t *testing.T) {
		_testSetContactAsPrimary(t, userService, userServiceText_ContactRepo)
	})

	t.Run("ConfirmContact", func(t *testing.T) {
		_testConfirmContact(t, userService, userServiceText_ContactRepo, userServiceText_TokenRepo)
	})

	t.Run("GetContactByID", func(t *testing.T) {
		_testGetContactByID(t, userService)
	})
}

func setupTestUserServiceData(t *testing.T, userRepo repo.UserRepo, contactRepo repo.ContactRepo) {
	// add user with confirmed contact
	userServiceTest_ConfirmedUser = models.User{
		PasswordHash: "does not matter",
	}
	err := userRepo.AddUser(context.TODO(), &userServiceTest_ConfirmedUser, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed primary contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail, core.Email, true)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedPrimaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmedprimary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed seconday contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail, core.Email, false)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user unconfirmed seconday contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail, core.Email, false)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed primary mobile contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedPrimaryMobile, core.Mobile, true)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed mobile primary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user unconfirmed secondary mobile contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobile, core.Mobile, false)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create unconfirmed mobile secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add unconfirmed user
	userServiceTest_UnconfirmedUser = models.User{
		PasswordHash: "does not matter",
	}
	err = userRepo.AddUser(context.TODO(), &userServiceTest_UnconfirmedUser, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create user with unconfirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add unconfrimed user unconfirmed primary contact
	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact = models.NewContact(userServiceTest_UnconfirmedUser.ID, "", userServiceTest_UnconfirmedUser_UnconfirmedPrimaryEmail, core.Email, true)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create unconfirmed primary contact for user with no confirmed contact for tests: %s", err.GetErrorCode())
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
	userServiceText_ContactRepo, err = memory.NewMemoryContactRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	userServiceText_TokenRepo = memory.NewMemoryTokenRepo()
	tokenService := NewTokenService(userServiceText_TokenRepo)
	userServiceTest_EmailService, _ = NewEmailService(StackEmailService, nil)
	userService := NewUserService(userRepo, userServiceText_ContactRepo, tokenService, userServiceTest_EmailService)
	setupTestUserServiceData(t, userRepo, userServiceText_ContactRepo)
	return userService
}

func _testUserServiceGetName(t *testing.T, userService services.UserService) {
	serviceName := userService.GetName()
	expectedServiceName := "userService"
	if serviceName != expectedServiceName {
		t.Errorf("\tservice name is not what was expected: got %s - expected %s", serviceName, expectedServiceName)
	}
}

func _testUserServiceGetUser(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		userID            string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:   "GIVEN a valid user id EXPECT user to be returned",
			userID: userServiceTest_ConfirmedUser.ID,
		},
		{
			name:              "GIVEN a non existant user id EXPECT error code no user found",
			userID:            "not-an-id", // TODO: need to have access to an id generator, so we can let this test work properly reguardless of the underlying data store used...
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := userService.GetUser(context.TODO(), logger, tc.userID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if user.ID != tc.userID {
					t.Errorf("\tuser id did not match expected: got - %s expected - %s", user.ID, tc.userID)
				}
			}
		})
	}
}

func _testGetUserAndContactByConfirmedContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		contactPrincipal  string
		contactType       core.ContactType
		expectedUserID    string
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "Given confirmed primary contact Return User",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail,
			contactType:       core.Email,
			expectedUserID:    userServiceTest_ConfirmedUser.ID,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
		},
		{
			name:              "Given confirmed secondary contact Return User",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail,
			contactType:       core.Email,
			expectedUserID:    userServiceTest_ConfirmedUser.ID,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
		},
		{
			name:              "Given unconfirmed contact Return error code no user found",
			contactPrincipal:  userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail,
			contactType:       core.Email,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
		{
			name:              "Given non existant contact Return No User Found Error",
			contactPrincipal:  "ojhgfiujwsfiogh@oiwujhgfiwsrb.dofuhsfoiuds",
			contactType:       core.Email,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, contact, err := userService.GetUserAndContactByConfirmedContact(context.TODO(), logger, tc.contactType, tc.contactPrincipal, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if user.ID != tc.expectedUserID {
					t.Errorf("\tuser id did not match expected: got - %s expected - %s", user.ID, tc.expectedUserID)
				}
				if contact.ID != tc.expectedContactID {
					t.Errorf("\tcontact id did not match expected: got - %s expected - %s", contact.ID, tc.expectedContactID)
				}
			}
		})
	}
}

func _testRegisterUserAndPrimaryContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		contactPrincipal  string
		contactType       core.ContactType
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:             "GIVEN unregistered contact EXPECT Successful registration new user and contact",
			contactPrincipal: userServiceTest_UserToRegisterEmail,
			contactType:      core.Email,
		},
		{
			name:             "GIVEN previously registered unconfirmed contact EXPECT successful registration new user and contact",
			contactPrincipal: userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail,
			contactType:      core.Email,
		},
		{
			name:              "GIVEN the provided contact is already confirmed in the data store EXPECT error contact already confirmed",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail,
			contactType:       core.Email,
			expectedErrorCode: coreerrors.ErrCodeRegistrationContactAlreadyConfirmed,
		},
		{
			name:              "GIVEN an invalid contact principal of the given type contact EXPECT error code invalid contact principal",
			contactPrincipal:  "not a valid email",
			contactType:       core.Email,
			expectedErrorCode: coreerrors.ErrCodeInvalidContactPrincipal,
		},
		{
			name:              "GIVEN an invalid contact principal of the given type contact EXPECT error code invalid contact principal",
			contactPrincipal:  "a_valid_email@email.com",
			contactType:       "not a valid contact type",
			expectedErrorCode: coreerrors.ErrCodeInvalidContactType,
		},
		// TODO: create test case for multiple confirmed instances of a contact returning the appropriate error...
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.RegisterUserAndPrimaryContact(context.TODO(), logger, tc.contactType, tc.contactPrincipal, "service_name", userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				ses, ok := userServiceTest_EmailService.(*stackEmailService)
				if !ok {
					t.Errorf("\texpected stackEmailService instance of email service but got %s", userServiceTest_EmailService.GetName())
					t.FailNow()
				}
				lastMessage, ok := ses.PopMessage()
				if !ok {
					t.Error("\tno message found in email stack from user registration")
				}
				numReceipents := len(lastMessage.To)
				if numReceipents != 1 {
					t.Errorf("\twrong number of recepitents in messages to: got - %d expected 1", numReceipents)
					// we return here because we do not need to go further
					return
				}
				if lastMessage.To[0] != tc.contactPrincipal {
					t.Errorf("\tto value not expected: got - %s expected - %s", lastMessage.To[0], tc.contactPrincipal)
				}
				// TODO: check subject and body once I write those....
			}
		})
	}
}

func _testGetUserPrimaryContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		userID            string
		contactType       core.ContactType
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN a valid user id EXPECT that users primary contact",
			userID:            userServiceTest_ConfirmedUser.ID,
			contactType:       core.Email,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := userService.GetUserPrimaryContact(context.TODO(), logger, tc.userID, tc.contactType, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if contact.ID != tc.expectedContactID {
					t.Errorf("\tcontact id not expected value: got - %s expected - %s", contact.ID, tc.expectedContactID)
				}
			}
		})
	}
}

func _testGetUsersContacts(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name               string
		userID             string
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:   "GIVEN a valid user id EXPECT all contacts associated with that user",
			userID: userServiceTest_ConfirmedUser.ID,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
				userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ID,
				userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact.ID,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contacts, err := userService.GetUsersContacts(context.TODO(), logger, tc.userID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
				} else {
					numContactsReturned := len(contacts)
					numExpectedContacts := len(tc.expectedContactIDs)
					if numContactsReturned != numExpectedContacts {
						t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
					}
					for _, ecid := range tc.expectedContactIDs {
						found := false
						for _, c := range contacts {
							if ecid == c.ID {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("\tunable to find expected contact id in results: %s", ecid)
						}
					}
				}
			}
		})
	}
}

func _testGetUsersConfirmedContacts(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name               string
		userID             string
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:   "GIVEN a valid user id EXPECT all confirmed contacts associated with that user",
			userID: userServiceTest_ConfirmedUser.ID,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ID,
			},
		},
		{
			name:               "GIVEN a valid user id with no confirmed contacts EXPECT no contacts to be returned",
			userID:             userServiceTest_UnconfirmedUser.ID,
			expectedContactIDs: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contacts, err := userService.GetUsersConfirmedContacts(context.TODO(), logger, tc.userID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
				} else {
					for _, ecid := range tc.expectedContactIDs {
						found := false
						for _, c := range contacts {
							if ecid == c.ID {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("\tunable to find expected contact id in results: %s", ecid)
						}
					}
				}
			}
		})
	}
}

func _testGetUsersContactsOfType(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name               string
		userID             string
		contactType        core.ContactType
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.Email,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
				userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user #2",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.Mobile,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ID,
				userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact.ID,
			},
		},
		{
			name:               "GIVEN a valid user id with no confirmed contacts EXPECT no contacts to be returned",
			userID:             userServiceTest_UnconfirmedUser.ID,
			expectedContactIDs: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contacts, err := userService.GetUsersContactsOfType(context.TODO(), logger, tc.userID, tc.contactType, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
				} else {
					for _, ecid := range tc.expectedContactIDs {
						found := false
						var contact models.Contact
						for _, c := range contacts {
							if ecid == c.ID {
								contact = c
								found = true
								break
							}
						}
						if !found {
							t.Errorf("\tunable to find expected contact id in results: %s", ecid)
						} else if contact.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", contact.Type, tc.contactType)
						}
					}
				}
			}
		})
	}
}

func _testGetUsersConfirmedContactsOfType(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name               string
		userID             string
		contactType        core.ContactType
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.Email,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user #2",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.Mobile,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ID,
			},
		},
		{
			name:               "GIVEN a valid user id with no confirmed contacts EXPECT no contacts to be returned",
			userID:             userServiceTest_UnconfirmedUser.ID,
			expectedContactIDs: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contacts, err := userService.GetUsersConfirmedContactsOfType(context.TODO(), logger, tc.userID, tc.contactType, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
				} else {
					for _, ecid := range tc.expectedContactIDs {
						found := false
						var contact models.Contact
						for _, c := range contacts {
							if ecid == c.ID {
								contact = c
								found = true
								break
							}
						}
						if !found {
							t.Errorf("\tunable to find expected contact id in results: %s", ecid)
						} else if contact.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", contact.Type, tc.contactType)
						}
					}
				}
			}
		})
	}
}

func _testAddContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		userID            string
		contactUserID     string
		contactPrincipal  string
		contactType       core.ContactType
		contactIsPrimary  bool
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:             "GIVEN a valid contact to add EXPECT the contact to be added",
			userID:           userServiceTest_ConfirmedUser.ID,
			contactUserID:    userServiceTest_ConfirmedUser.ID,
			contactPrincipal: "contact_to_add_7y6egdf7cya@email.com",
			contactType:      core.Email,
			contactIsPrimary: false,
		},
		{
			name:             "GIVEN a contact to add that already exists elsewhere but is not confirmed EXPECT the contact to be added",
			userID:           userServiceTest_UnconfirmedUser.ID,
			contactUserID:    userServiceTest_UnconfirmedUser.ID,
			contactPrincipal: userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.Principal,
			contactType:      userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.Type,
			contactIsPrimary: false,
		},
		{
			name:              "GIVEN a valid contact to add that is marked as primary EXPECT error code contact to add marked as primary",
			userID:            userServiceTest_ConfirmedUser.ID,
			contactUserID:     userServiceTest_ConfirmedUser.ID,
			contactPrincipal:  "contact_to_add_208yu7fiwr@email.com",
			contactType:       core.Email,
			contactIsPrimary:  true,
			expectedErrorCode: coreerrors.ErrCodeContactToAddMarkedAsPrimary,
		},
		{
			name:              "GIVEN a contact that is already confirmed somewhere else EXPECT error code contact to add already confirmed",
			userID:            userServiceTest_UnconfirmedUser.ID,
			contactUserID:     userServiceTest_UnconfirmedUser.ID,
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.Principal,
			contactType:       userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.Type,
			contactIsPrimary:  false,
			expectedErrorCode: coreerrors.ErrCodeContactToAddAlreadyConfirmed,
		},
		{
			name:              "GIVEN a valid contact to add and a non matching user id EXPECT error code user ids do not match",
			userID:            userServiceTest_UnconfirmedUser.ID,
			contactUserID:     userServiceTest_ConfirmedUser.ID,
			contactPrincipal:  "contact_to_add_56yhgvfdew@email.com",
			contactType:       core.Email,
			contactIsPrimary:  false,
			expectedErrorCode: coreerrors.ErrCodeUserIDsDoNotMatch,
		},
		{
			name:              "GIVEN an invalid contact principal of the given type contact EXPECT error code invalid contact principal",
			userID:            userServiceTest_UnconfirmedUser.ID,
			contactUserID:     userServiceTest_UnconfirmedUser.ID,
			contactPrincipal:  "not a valid email",
			contactType:       core.Email,
			contactIsPrimary:  false,
			expectedErrorCode: coreerrors.ErrCodeInvalidContactPrincipal,
		},
		{
			name:              "GIVEN an invalid contact type of the given type contact EXPECT error code invalid contact type",
			userID:            userServiceTest_UnconfirmedUser.ID,
			contactUserID:     userServiceTest_UnconfirmedUser.ID,
			contactPrincipal:  "a_valid_email@email.com",
			contactType:       "not a valid contact type",
			contactIsPrimary:  false,
			expectedErrorCode: coreerrors.ErrCodeInvalidContactType,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newContact := models.NewContact(tc.contactUserID, "", tc.contactPrincipal, tc.contactType, tc.contactIsPrimary)
			err := userService.AddContact(context.TODO(), logger, tc.userID, &newContact, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if newContact.ID == "" {
					t.Error("\tcontact added does not have a contact id")
				}
				if newContact.UserID != tc.userID {
					t.Errorf("\tadded contact user id does not match expected user id: got - %s expected - %s", newContact.UserID, tc.userID)
				}
				// test that confirmation email was sent
				ses, ok := userServiceTest_EmailService.(*stackEmailService)
				if !ok {
					t.Errorf("\texpected stackEmailService instance of email service but got %s", userServiceTest_EmailService.GetName())
					t.FailNow()
				}
				lastMessage, ok := ses.PopMessage()
				if !ok {
					t.Error("\tno message found in email stack from user registration")
				}
				numReceipents := len(lastMessage.To)
				if numReceipents != 1 {
					t.Errorf("\twrong number of recepitents in messages to: got - %d expected 1", numReceipents)
					// we return here because we do not need to go further
					return
				}
				if lastMessage.To[0] != tc.contactPrincipal {
					t.Errorf("\tto value not expected: got - %s expected - %s", lastMessage.To[0], tc.contactPrincipal)
				}
				// TODO: check subject and body once I write those....
			}
		})
	}
}

func _testSetContactAsPrimary(t *testing.T, userService services.UserService, contactRepo repo.ContactRepo) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name                            string
		userID                          string
		contactType                     core.ContactType
		newPrimaryContactID             string
		expectedCurrentPrimaryContactID string
		expectedErrorCode               string
	}
	testCases := []testCase{
		{
			name:                            "GIVEN a proper contact id EXPECT new primary contact to be set as primary, and the old primary to be set as not primary",
			userID:                          userServiceTest_ConfirmedUser.ID,
			contactType:                     userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.Type,
			newPrimaryContactID:             userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
			expectedCurrentPrimaryContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
		},
		{
			name:                            "GIVEN a proper contact id EXPECT new primary contact to be set as primary, and the old primary to be set as not primary (revert previous)",
			userID:                          userServiceTest_ConfirmedUser.ID,
			contactType:                     userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.Type,
			newPrimaryContactID:             userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedCurrentPrimaryContactID: userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
		},
		{
			name:                            "GIVEN a contact id that is not confirmed EXPECT error code contact not confirmed",
			userID:                          userServiceTest_ConfirmedUser.ID,
			contactType:                     userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.Type,
			newPrimaryContactID:             userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.ID,
			expectedCurrentPrimaryContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedErrorCode:               coreerrors.ErrCodeContactNotConfirmed,
		},
		{
			name:                            "GIVEN a contact id is already marked as primary EXPECT error code contact already marked as primary",
			userID:                          userServiceTest_ConfirmedUser.ID,
			contactType:                     userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.Type,
			newPrimaryContactID:             userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedCurrentPrimaryContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedErrorCode:               coreerrors.ErrCodeContactAlreadyMarkedPrimary,
		},
		{
			name:                            "GIVEN a contact id thats associated user id does not match EXPECT error code user ids do not match",
			userID:                          userServiceTest_UnconfirmedUser.ID,
			contactType:                     userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.Type,
			newPrimaryContactID:             userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
			expectedCurrentPrimaryContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedErrorCode:               coreerrors.ErrCodeUserIDsDoNotMatch,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentPrimaryContact, err := contactRepo.GetPrimaryContactByUserID(context.TODO(), tc.userID, tc.contactType)
			if err != nil {
				t.Errorf("\tfailed to retreive current primary contact for user")
				return
			}
			err = userService.SetContactAsPrimary(context.TODO(), logger, tc.userID, tc.newPrimaryContactID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				contactsOfType, err := contactRepo.GetContactsByUserIDAndType(context.TODO(), tc.userID, tc.contactType)
				if err != nil {
					t.Errorf("\tfailed to return users contacts for test confirmation: %s - %s", err.GetErrorCode(), err.Error())
					return
				}
				var previousPrimaryContact, newPrimaryContact models.Contact
				var numPrimaryContacts = 0
				for _, c := range contactsOfType {
					if c.ID == currentPrimaryContact.ID {
						previousPrimaryContact = c
					} else if c.ID == tc.newPrimaryContactID {
						newPrimaryContact = c
					}
					if c.IsPrimary {
						numPrimaryContacts++
					}

				}
				if numPrimaryContacts != 1 {
					t.Errorf("\tthere should only be one primary contact, but there were %d", numPrimaryContacts)
				}
				if previousPrimaryContact.IsPrimary {
					t.Error("\tprevious primary contact should not still be marked as primary")
				}
				if !newPrimaryContact.IsPrimary {
					t.Error("\tnew primary contact should be marked as primary")
				}
				if currentPrimaryContact.ID != tc.expectedCurrentPrimaryContactID {
					t.Errorf("\tprevious contact id not expected value: got - %s expected - %s", currentPrimaryContact.ID, tc.expectedCurrentPrimaryContactID)
				}
				if newPrimaryContact.ID != tc.newPrimaryContactID {
					t.Errorf("\tprevious contact id not expected value: got - %s expected - %s", newPrimaryContact.ID, tc.newPrimaryContactID)
				}
			}
		})
	}
}

func _testConfirmContact(t *testing.T, userService services.UserService, contactRepo repo.ContactRepo, tokenRepo repo.TokenRepo) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name                          string
		contactToConfirm              *models.Contact
		tokenValidFor                 time.Duration
		mockConfirmContactTokenString string
		expectedErrorCode             string
	}
	testCases := []testCase{
		{
			name:                          "GIVEN a valid confirmation code for a contact to be confirmed EXPECT success and contact confirmation date to be set",
			contactToConfirm:              &userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact,
			tokenValidFor:                 time.Minute,
			mockConfirmContactTokenString: "",
			expectedErrorCode:             "",
		},
		{
			name:                          "GIVEN an invalid confirmation code for a contact to be confirmed EXPECT error code invalid token",
			contactToConfirm:              &userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact,
			tokenValidFor:                 time.Minute,
			mockConfirmContactTokenString: "whyuibgvouieynboiuwyb04t8b5tu7yv394uy9tur",
			expectedErrorCode:             coreerrors.ErrCodeInvalidToken,
		},
		{
			name:                          "GIVEN an expired confirmation code for a contact to be confirmed EXPECT error code invalid token",
			contactToConfirm:              &userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact,
			tokenValidFor:                 time.Minute * -1,
			mockConfirmContactTokenString: "",
			expectedErrorCode:             coreerrors.ErrCodeExpiredToken,
		},
		{
			name:                          "GIVEN an expired confirmation code for a contact to be confirmed EXPECT error code invalid token",
			contactToConfirm:              &userServiceTest_ConfirmedUser_ConfirmedSecondaryContact,
			tokenValidFor:                 time.Minute,
			mockConfirmContactTokenString: "",
			expectedErrorCode:             coreerrors.ErrCodeContactAlreadyConfirmed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var confirmContactToken string
			if tc.mockConfirmContactTokenString != "" {
				confirmContactToken = tc.mockConfirmContactTokenString
			} else {
				newConfirmToken, err := models.NewToken(tc.contactToConfirm.ID, models.TokenTypeConfirmContact, tc.tokenValidFor)
				if err != nil {
					t.Errorf("\tfailed to create new confirm contact token: %s - %s", err.GetErrorCode(), err.Error())
				}
				err = tokenRepo.PutToken(context.TODO(), newConfirmToken)
				if err != nil {
					t.Errorf("\tfailed to new confirm contact token in repo for validation: %s - %s", err.GetErrorCode(), err.Error())
				}
				confirmContactToken = newConfirmToken.Value
			}
			err := userService.ConfirmContact(context.TODO(), logger, confirmContactToken, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				newlyConfirmedContact, err := contactRepo.GetContactByID(context.TODO(), tc.contactToConfirm.ID)
				if err != nil {
					t.Errorf("\tfailed to retreive newly confirmed contact from repo for validation: %s - %s", err.GetErrorCode(), err.Error())
				}
				if !newlyConfirmedContact.IsConfirmed() {
					t.Error("\tnewly confirmed contact is not confirmed in the underlying data store.")
				}
			}
		})
	}
}

func _testGetContactByID(t *testing.T, userService services.UserService) {
	type testCase struct {
		name                     string
		contactID                string
		expectedContactType      core.ContactType
		expectedContactPrincipal string
		expectedErrorCode        string
	}
	testCases := []testCase{
		{
			name:                     "GIVEN an existing contact id EXPECT success",
			contactID:                userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
			expectedContactPrincipal: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.Principal,
			expectedContactType:      userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.Type,
		},
		{
			name:              "GIVEN a non existant contact id EXPECT error code no contact found",
			contactID:         "not a real id",
			expectedErrorCode: coreerrors.ErrCodeNoContactFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			contact, err := userService.GetContactByID(context.TODO(), logger, tc.contactID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if contact.Principal != tc.expectedContactPrincipal {
					t.Errorf("\tPrincipal value was not expected: got - %s expected - %s", contact.Principal, tc.expectedContactPrincipal)
				}
				if contact.Type != tc.expectedContactType {
					t.Errorf("\tType value was not expected: got - %s expected - %s", contact.Type, tc.expectedContactType)
				}
			}
		})
	}
}
