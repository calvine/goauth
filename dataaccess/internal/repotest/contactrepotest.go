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
)

const (
	contactRepoCreatedBy = "contact repo tests"
)

var (
	newContact1           models.Contact
	newContact2           models.Contact
	newContact3           models.Contact
	newContact4           models.Contact
	noMatchingUserContact models.Contact

	nonExistantContactID string
)

// func TestInterface(t *testing.T) {
// 	cr := NewUserRepo(nil)
// 	reflect.ValueOf(cr).
// }

func setupContactTestData(t *testing.T, testHarness RepoTestHarnessInput) {
	newContact1 = models.Contact{
		Principal:     "testuser1@domain.org",
		UserID:        testUser1.ID,
		IsPrimary:     true,
		Type:          core.CONTACT_TYPE_EMAIL,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newContact2 = models.Contact{
		Principal:     "frank@app.space",
		UserID:        testUser1.ID,
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_EMAIL,
		ConfirmedDate: nullable.NullableTime{},
	}

	newContact3 = models.Contact{
		Principal:     "555-555-5555",
		UserID:        testUser1.ID,
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_MOBILE,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newContact4 = models.Contact{
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
	// t.Run("GetContactByConfirmationCode", func(t *testing.T) {
	// 	_testGetContactByConfirmationCode(t, contactRepo)
	// })
	t.Run("UpdateContact", func(t *testing.T) {
		_testUpdateContact(t, *testHarness.ContactRepo)
	})
	// t.Run("ConfirmContact", func(t *testing.T) {
	// 	_testConfirmContact(t, contactRepo)
	// })
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
			contactToAdd:   &newContact1,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newContact2,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newContact3,
			expectedUserID: testUser1.ID,
		},
		{
			name:           "GIVEN contact to insert EXPECT contact to be inserted",
			contactToAdd:   &newContact4,
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
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
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
			contactID: newContact1.ID,
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
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
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
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN an existing user id with a primary contact EXPECT the primary contact to be returned",
			userID:            testUser1.ID,
			expectedContactID: newContact1.ID,
		},
		{
			name:              "GIVEN an nonexistant user id EXPECT error no user found",
			userID:            nonExistantUserID,
			expectedErrorCode: coreerrors.ErrCodeNoUserFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := contactRepo.GetPrimaryContactByUserID(context.TODO(), tc.userID)
			if err != nil {
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
			} else {
				if !contact.IsPrimary {
					t.Errorf("\t contact returned is not primary contact: %v", contact)
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
				newContact1.ID,
				newContact2.ID,
				newContact3.ID,
				newContact4.ID,
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
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
			} else {
				numContactsFound := len(contacts)
				if numContactsFound != numExpectedContacts {
					t.Errorf("\tnumber of contacts returned not expected amount: got - %d expected - %d", numContactsFound, numExpectedContacts)
					t.Fail()
				} else {
					numPrimary := 0
					for _, c := range contacts {
						contactFound := false
						if c.IsPrimary {
							numPrimary++
						}
						for _, ecid := range tc.expectedContactIds {
							if c.ID == ecid {
								contactFound = true
								break
							}
						}
						if !contactFound {
							t.Errorf("\tcontact with id %swas not found in results", c.ID)
							t.Fail()
						}
					}
					if numPrimary != 1 {
						t.Errorf("\texpected 1 primary contact but got %d", numPrimary)
					}
				}
			}
		})
	}
}

// func _testGetContactByConfirmationCode(t *testing.T, userRepo repo.ContactRepo) {
// 	confirmationCode := newContact2.ConfirmationCode.Value
// 	expectedPrincipal := newContact2.Principal
// 	expectedIsPrimary := newContact2.IsPrimary
// 	contact, err := userRepo.GetContactByConfirmationCode(context.TODO(), confirmationCode)
// 	if err != nil {
// 		t.Error("failed to get contact by confirmation code", err)
// 	}
// 	if contact.Principal != expectedPrincipal {
// 		t.Error("unexpected value of Principal", expectedPrincipal, contact.Principal)
// 	}
// 	if contact.IsPrimary != expectedIsPrimary {
// 		t.Error("unexpected value of IsPrimary", expectedIsPrimary, contact.IsPrimary)
// 	}
// 	if contact.ConfirmationCode.Value != confirmationCode {
// 		t.Error("expected retreived confirmation code to match passed in confirmation code.", confirmationCode, contact.ConfirmationCode.Value)
// 	}
// }

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
			contactToUpdate: &newContact3,
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
				if tc.expectedErrorCode == "" {
					t.Errorf("\tunexpected error encountered: %s - %s", err.GetErrorCode(), err.Error())
					t.Fail()
				} else if tc.expectedErrorCode != err.GetErrorCode() {
					t.Errorf("\terror code did not match expected: got - %s expected - %s", err.GetErrorCode(), tc.expectedErrorCode)
					t.Fail()
				}
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
