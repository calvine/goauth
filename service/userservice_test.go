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

	// t.Run("SetContactAsPrimary", func(t *testing.T) {
	// 	_testSetContactAsPrimary(t, userService)
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
		t.Errorf("\tfailed to create user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed primary contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, true)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedPrimaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmedprimary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed seconday contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, false)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user unconfirmed seconday contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, false)
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed secondary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user confirmed primary mobile contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_ConfirmedPrimaryMobile, core.CONTACT_TYPE_MOBILE, true)
	// confirm the contact
	userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact.ConfirmedDate.Set(time.Now().Add(time.Minute * -1))
	err = contactRepo.AddContact(context.TODO(), &userServiceTest_ConfirmedUser_ConfirmedPrimaryMobileContact, userServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("\tfailed to create confirmed mobile primary contact of user with confirmed contact for tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	// add confirmed user unconfirmed secondary mobile contact
	userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobileContact = models.NewContact(userServiceTest_ConfirmedUser.ID, "", userServiceTest_ConfirmedUser_UnconfirmedSecondaryMobile, core.CONTACT_TYPE_MOBILE, false)
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
	userServiceTest_UnconfirmedUser_UnconfirmedPrimaryContact = models.NewContact(userServiceTest_UnconfirmedUser.ID, "", userServiceTest_UnconfirmedUser_UnconfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, true)
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
	contactRepo, err := memory.NewMemoryContactRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tokenRepo := memory.NewMemoryTokenRepo()
	tokenService := NewTokenService(tokenRepo)
	userServiceTest_EmailService, _ = NewEmailService(StackEmailService, nil)
	userService := NewUserService(userRepo, contactRepo, tokenService, userServiceTest_EmailService)
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
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
			name:             "GIVEN unregistered contact EXPECT Successful registration new user and contact",
			contactPrincipal: userServiceTest_UserToRegisterEmail,
			contactType:      core.CONTACT_TYPE_EMAIL,
		},
		{
			name:             "GIVEN previously registered unconfirmed contact EXPECT successful registration new user and contact",
			contactPrincipal: userServiceTest_ConfirmedUser_UnconfirmedSecondaryEmail,
			contactType:      core.CONTACT_TYPE_EMAIL,
		},
		{
			name:              "GIVEN the provided contact is already confirmed in the data store EXPECT error contact already confirmed",
			contactPrincipal:  userServiceTest_ConfirmedUser_ConfirmedSecondaryEmail,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeRegistrationContactAlreadyConfirmed,
		},
		// TODO: create test case for multiple confirmed instances of a contact returning the appropriate error...
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.RegisterUserAndPrimaryContact(context.TODO(), logger, tc.contactType, tc.contactPrincipal, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				ses, ok := userServiceTest_EmailService.(*stackEmailService)
				if !ok {
					t.Errorf("\texpected stackEmailService instance of email service but got %s", userServiceTest_EmailService.GetName())
					t.FailNow()
				}
				lastMessage, ok := ses.PopMessage()
				if !ok {
					t.Error("\tno message found in email stack from user registration")
					t.Fail()
				}
				numReceipents := len(lastMessage.To)
				if numReceipents != 1 {
					t.Errorf("\twrong number of recepitents in messages to: got - %d expected 1", numReceipents)
					t.Fail()
					// we return here because we do not need to go further
					return
				}
				if lastMessage.To[0] != tc.contactPrincipal {
					t.Errorf("\tto value not expected: got - %s expected - %s", lastMessage.To[0], tc.contactPrincipal)
					t.Fail()
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
		contactType       string
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN a valid user id EXPECT that users primary contact",
			userID:            userServiceTest_ConfirmedUser.ID,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedContactID: userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := userService.GetUserPrimaryContact(context.TODO(), logger, tc.userID, tc.contactType, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				if contact.ID != tc.expectedContactID {
					t.Errorf("\tcontact id not expected value: got - %s expected - %s", contact.ID, tc.expectedContactID)
					t.Fail()
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
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
					t.Fail()
				} else {
					numContactsReturned := len(contacts)
					numExpectedContacts := len(tc.expectedContactIDs)
					if numContactsReturned != numExpectedContacts {
						t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
						t.Fail()
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
							t.Fail()
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
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
					t.Fail()
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
							t.Fail()
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
		contactType        string
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.CONTACT_TYPE_EMAIL,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
				userServiceTest_ConfirmedUser_UnconfirmedSecondaryContact.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user #2",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.CONTACT_TYPE_MOBILE,
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
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
					t.Fail()
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
							t.Fail()
						} else if contact.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", contact.Type, tc.contactType)
							t.Fail()
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
		contactType        string
		expectedContactIDs []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.CONTACT_TYPE_EMAIL,
			expectedContactIDs: []string{
				userServiceTest_ConfirmedUser_ConfirmedPrimaryContact.ID,
				userServiceTest_ConfirmedUser_ConfirmedSecondaryContact.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type EXPECT all confirmed contacts associated with that user #2",
			userID:      userServiceTest_ConfirmedUser.ID,
			contactType: core.CONTACT_TYPE_MOBILE,
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
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				numContactsReturned := len(contacts)
				numExpectedContacts := len(tc.expectedContactIDs)
				if numContactsReturned != numExpectedContacts {
					t.Errorf("\tnumber of contacts expected does not match: got - %d expected - %d", numContactsReturned, numExpectedContacts)
					t.Fail()
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
							t.Fail()
						} else if contact.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", contact.Type, tc.contactType)
							t.Fail()
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
		contactType       string
		contactIsPrimary  bool
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:             "GIVEN a valid contact to add EXPECT the contact to be added",
			userID:           userServiceTest_ConfirmedUser.ID,
			contactUserID:    userServiceTest_ConfirmedUser.ID,
			contactPrincipal: "contact_to_add_7y6egdf7cya@email.com",
			contactType:      core.CONTACT_TYPE_EMAIL,
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
			contactType:       core.CONTACT_TYPE_EMAIL,
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
			contactType:       core.CONTACT_TYPE_EMAIL,
			contactIsPrimary:  false,
			expectedErrorCode: coreerrors.ErrCodeUserIDsDoNotMatch,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newContact := models.NewContact(tc.contactUserID, "", tc.contactPrincipal, tc.contactType, tc.contactIsPrimary)
			err := userService.AddContact(context.TODO(), logger, tc.userID, &newContact, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				if newContact.ID == "" {
					t.Error("\tcontact added does not have a contact id")
					t.Fail()
				}
				if newContact.UserID != tc.userID {
					t.Errorf("\tadded contact user id does not match expected user id: got - %s expected - %s", newContact.UserID, tc.userID)
					t.Fail()
				}
			}
		})
	}
}

func _testSetContactAsPrimary(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name                string
		userID              string
		newPrimaryContactID string
		expectedErrorCode   string
	}
	testCases := []testCase{
		{
			name: "GIVEN a proper contact id EXPECT new primary contact to be set as primary, and the old primary to be set as not primary",
		},
		{
			name: "GIVEN a proper contact id EXPECT new primary contact to be set as primary, and the old primary to be set as not primary",
		},
		{
			name: "GIVEN a contact id that is not confirmed EXPECT error code contact not confirmed",
		},
		{
			name: "GIVEN a contact id is already marked as primary EXPECT error code contact already marked as primary",
		},
		{
			name: "GIVEN a contact id thats associated user id does not match EXPECT error code user ids do not match",
		},
		{
			name: "GIVEN a contact id that is not confirmed EXPECT error code contact not confirmed",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.SetContactAsPrimary(context.TODO(), logger, tc.userID, tc.newPrimaryContactID, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				// TODO: check that data store state is correct (only one primary contact)
			}
		})
	}
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}

func _testConfirmContact(t *testing.T, userService services.UserService) {
	logger := zaptest.NewLogger(t)
	type testCase struct {
		name              string
		confirmationCode  string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid confirmation code for a contact to be confirmed EXPECT success and contact confirmation date to be set",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userService.ConfirmContact(context.TODO(), logger, tc.confirmationCode, userServiceTest_CreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("expected an error to occurr: %s", tc.expectedErrorCode)
				t.Fail()
			} else {
				// TODO: check that data store contact is confirmed
			}
			// There is nothing to check really?
		})
	}
	t.Error(coreerrors.NewNotImplementedError(true))
	t.Fail()
}
