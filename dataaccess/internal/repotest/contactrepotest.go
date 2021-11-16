package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

const (
	contactRepoCreatedBy = "contact repo tests"
)

var (
	newEmailContact1      models.Contact
	newEmailContact2      models.Contact
	newMobileContact1     models.Contact
	newMobileContact2     models.Contact
	noMatchingUserContact models.Contact

	nonExistantContactID string
)

// func TestInterface(t *testing.T) {
// 	cr := NewUserRepo(nil)
// 	reflect.ValueOf(cr).
// }

func setupContactTestData(t *testing.T, testHarness RepoTestHarnessInput) {
	newEmailContact1 = models.Contact{
		Principal:     "testuser1@domain.org",
		UserID:        testUser1.ID,
		IsPrimary:     true,
		Type:          core.CONTACT_TYPE_EMAIL,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newEmailContact2 = models.Contact{
		Principal:     "frank@app.space",
		UserID:        testUser1.ID,
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_EMAIL,
		ConfirmedDate: nullable.NullableTime{},
	}

	newMobileContact1 = models.Contact{
		Principal:     "555-555-5555",
		UserID:        testUser1.ID,
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_MOBILE,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newMobileContact2 = models.Contact{
		Principal:     "email@email.email",
		UserID:        testUser1.ID,
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_MOBILE,
		ConfirmedDate: nullable.NullableTime{HasValue: false},
	}

	noMatchingUserContact = models.Contact{
		Principal:     "email@email.email",
		UserID:        testHarness.IDGenerator(false),
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_MOBILE,
		ConfirmedDate: nullable.NullableTime{HasValue: false},
	}

	nonExistantContactID = testHarness.IDGenerator(false)
}

func testContactRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	setupContactTestData(t, testHarness)
	t.Run("AddContact", func(t *testing.T) {
		_testAddContact(t, *testHarness.ContactRepo)
	})
	t.Run("GetContactByID", func(t *testing.T) {
		_testGetContactByID(t, *testHarness.ContactRepo)
	})
	t.Run("GetPrimaryContactByUserID", func(t *testing.T) {
		_testGetPrimaryContactByUserID(t, *testHarness.ContactRepo)
	})
	t.Run("GetContactsByUserID", func(t *testing.T) {
		_testGetContactsByUserID(t, *testHarness.ContactRepo)
	})
	t.Run("GetContactsByUserIDAndType", func(t *testing.T) {
		_testGetContactsByUserIDAndType(t, *testHarness.ContactRepo)
	})
	t.Run("UpdateContact", func(t *testing.T) {
		_testUpdateContact(t, *testHarness.ContactRepo)
	})
	t.Run("GetExistingConfirmedContactsCountByPrincipalAndType", func(t *testing.T) {
		_testGetExistingConfirmedContactsCountByPrincipalAndType(t, *testHarness.ContactRepo)
	})
}

func _testAddContact(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		contactToAdd      *models.Contact
		expectedUserID    string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newEmailContact1,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newEmailContact2,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newMobileContact1,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newMobileContact2,
			expectedUserID: testUser1.ID,
		},
		{
			name:              "GIVEN contact with non existant user id EXPECT no user found error code",
			contactToAdd:      &noMatchingUserContact,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := contactRepo.AddContact(context.TODO(), tc.contactToAdd, contactRepoCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if tc.expectedUserID != tc.contactToAdd.UserID {
					t.Errorf("\tuser id expected: got - %s expected - %s", tc.contactToAdd.UserID, tc.expectedUserID)
					t.Fail()
				}
				if tc.contactToAdd.ID == "" {
					t.Error("\tcontact id is blank")
					t.Fail()
				}
			}
		})
	}
}

func _testGetContactByID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		contactID         string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:      "GIVEN an existing contact id EXPECT the contact to be returned",
			contactID: newEmailContact1.ID,
		},
		{
			name:              "GIVEN a nonexistant contact id EXPECT error no contact found",
			contactID:         nonExistantContactID,
			expectedErrorCode: coreerrors.ErrCodeNoContactFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := contactRepo.GetContactByID(context.TODO(), tc.contactID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if tc.contactID != contact.ID {
					t.Errorf("\tcontact id expected: got - %s expected - %s", contact.ID, tc.contactID)
					t.Fail()
				}
				if contact.ID == "" {
					t.Error("\tcontact id is blank")
					t.Fail()
				}
			}
		})
	}
}

func _testGetPrimaryContactByUserID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		userID            string
		contactType       string
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN an existing user id with a primary contact EXPECT the primary contact to be returned",
			userID:            testUser1.ID,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedContactID: newEmailContact1.ID,
		},
		{
			name:              "GIVEN an nonexistant user id EXPECT error no user found",
			userID:            nonExistantUserID,
			contactType:       core.CONTACT_TYPE_EMAIL,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := contactRepo.GetPrimaryContactByUserID(context.TODO(), tc.userID, tc.contactType)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if !contact.IsPrimary {
					t.Errorf("\tcontact returned is not primary contact: %v", contact)
				}
				if tc.expectedContactID != contact.ID {
					t.Errorf("\tcontact id was not expected: got - %s expected - %s", contact.ID, tc.expectedContactID)
					t.Fail()
				}
				if contact.ID == "" {
					t.Error("\tcontact id is blank")
					t.Fail()
				}
			}
		})
	}
}

func _testGetContactsByUserID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name               string
		userID             string
		expectedContactIds []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:   "GIVEN a user id for a user with one or more contacts EXPECT the contacts for that user to be returned",
			userID: testUser1.ID,
			expectedContactIds: []string{
				newEmailContact1.ID,
				newEmailContact2.ID,
				newMobileContact1.ID,
				newMobileContact2.ID,
			},
		},
		{
			name:              "GIVEN a nonexistant user id EXPECT error no user found",
			userID:            nonExistantUserID,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			numExpectedContacts := len(tc.expectedContactIds)
			contacts, err := contactRepo.GetContactsByUserID(context.TODO(), tc.userID)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				numContactsFound := len(contacts)
				if numContactsFound != numExpectedContacts {
					t.Errorf("\tnumber of contacts returned not expected amount: got - %d expected - %d", numContactsFound, numExpectedContacts)
					t.Fail()
				} else {
					for _, c := range contacts {
						contactFound := false
						for _, ecid := range tc.expectedContactIds {
							if c.ID == ecid {
								contactFound = true
								break
							}
						}
						if !contactFound {
							t.Errorf("\tcontact with id %s was not found in results", c.ID)
							t.Fail()
						}
					}
				}
			}
		})
	}
}

func _testGetContactsByUserIDAndType(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name               string
		userID             string
		contactType        string
		expectedContactIds []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT contacts of that type returned",
			userID:      testUser1.ID,
			contactType: core.CONTACT_TYPE_EMAIL,
			expectedContactIds: []string{
				newEmailContact1.ID,
				newEmailContact2.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type #2 EXPECT contacts of that type returned",
			userID:      testUser1.ID,
			contactType: core.CONTACT_TYPE_MOBILE,
			expectedContactIds: []string{
				newMobileContact1.ID,
				newMobileContact2.ID,
			},
		},
		{
			name:               "GIVEN a valid user id and contact type that has no contacts of that type EXPECT 0 contacts to be returned",
			userID:             testUser1.ID,
			contactType:        "not a valid type", // TODO: need to make a case useing a valid contact type...
			expectedContactIds: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			numExpectedContacts := len(tc.expectedContactIds)
			contacts, err := contactRepo.GetContactsByUserIDAndType(context.TODO(), tc.userID, tc.contactType)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				numContactsFound := len(contacts)
				if numContactsFound != numExpectedContacts {
					t.Errorf("\tnumber of contacts returned not expected amount: got - %d expected - %d", numContactsFound, numExpectedContacts)
					t.Fail()
				} else {
					for _, c := range contacts {
						contactFound := false
						for _, ecid := range tc.expectedContactIds {
							if c.ID == ecid {
								contactFound = true
								break
							}
						}
						if !contactFound {
							t.Errorf("\tcontact with id %s was not found in results", c.ID)
							t.Fail()
						} else if c.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", c.Type, tc.contactType)
							t.Fail()
						}
					}
				}
			}
		})
	}
}

func _testUpdateContact(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		contactToUpdate   *models.Contact
		newPrincipal      string
		markConfirmed     bool
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:            "GIVEN a contact to update EXPECT contact to be updated",
			contactToUpdate: &newMobileContact1,
			newPrincipal:    "a_different_email@mail.org",
			markConfirmed:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			preUpdateTime := time.Now().UTC().Add(time.Second * -1)
			tc.contactToUpdate.Principal = tc.newPrincipal
			if tc.markConfirmed {
				tc.contactToUpdate.ConfirmedDate.Set(time.Now().UTC())
			}
			err := contactRepo.UpdateContact(context.TODO(), tc.contactToUpdate, contactRepoCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if tc.contactToUpdate.Principal != tc.newPrincipal {
					t.Errorf("\tupdated contact principal was not expected: got - %s expected - %s", tc.contactToUpdate.Principal, tc.newPrincipal)
					t.Fail()
				}
				if tc.markConfirmed {
					if !tc.contactToUpdate.ConfirmedDate.HasValue {
						t.Error("\tupdated contact confirmed date does not have a value even thought test case markConfirmed was true")
						t.Fail()
					}
					if tc.contactToUpdate.ConfirmedDate.Value.Before(preUpdateTime) {
						t.Errorf("\tupdated contact confirmed date not after pre update time when markConfirmed was set: preUpdateTime - %s confirmedDate - %s", preUpdateTime.Format(time.RFC3339), tc.contactToUpdate.ConfirmedDate.Value.Format(time.RFC3339))
						t.Fail()
					}
				}
				if tc.contactToUpdate.AuditData.ModifiedOnDate.HasValue &&
					tc.contactToUpdate.AuditData.ModifiedOnDate.Value.Before(preUpdateTime) {
					t.Error("\texpected updated contact ModifiedOnDate to be after the preUpdateTime")
					t.Fail()
				}
				if tc.contactToUpdate.AuditData.ModifiedByID.HasValue &&
					tc.contactToUpdate.AuditData.ModifiedByID.Value != contactRepoCreatedBy {
					t.Errorf("\tupdated contact ModifiedByID not expected: got %s - expected: %s", tc.contactToUpdate.AuditData.ModifiedByID.Value, contactRepoCreatedBy)
					t.Fail()
				}
			}
		})
	}
}

func _testGetExistingConfirmedContactsCountByPrincipalAndType(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name                     string
		contactPrincipal         string
		contactType              string
		expectedNumberOfContacts int64
		expectedErrorCode        string
	}
	testCases := []testCase{
		{
			name:                     "GIVEN a confirmed contact EXPECT 1 to be returned",
			contactPrincipal:         initialTestConfirmedPrimaryContact.Principal,
			contactType:              initialTestConfirmedPrimaryContact.Type,
			expectedNumberOfContacts: 1,
		},
		{
			name:                     "GIVEN an unconfirmed contact EXPECT 0 to be returned",
			contactPrincipal:         initialTestUnconfirmedContact.Principal,
			contactType:              initialTestUnconfirmedContact.Type,
			expectedNumberOfContacts: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			numContacts, err := contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType(context.TODO(), tc.contactType, tc.contactPrincipal)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else {
				if numContacts != tc.expectedNumberOfContacts {
					t.Errorf("\tnumber of contacts returned does not match expected: got - %d expected - %d", numContacts, tc.expectedNumberOfContacts)
					t.Fail()
				}
			}
		})
	}
}
