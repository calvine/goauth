package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/internal/testutils"
)

const (
	contactRepoCreatedBy = "contact repo tests"
)

var (
	newEmailContact1      models.Contact
	newEmailContact2      models.Contact
	newEmailContact3      models.Contact
	newMobileContact1     models.Contact
	newMobileContact2     models.Contact
	noMatchingUserContact models.Contact

	nonExistantContactID string
)

func setupContactTestData(t *testing.T, testHarness RepoTestHarnessInput) {
	newEmailContact1 = models.NewContact(testUser1.ID, "", "testuser1@domain.org", core.Email, true)
	newEmailContact1.ConfirmedDate.Set(time.Now().UTC())

	newEmailContact2 = models.NewContact(testUser1.ID, "", "frank@app.space", core.Email, false)

	newEmailContact3 = models.NewContact(testUser1.ID, "", "jeremy@app.space", core.Email, false)
	newEmailContact3.ConfirmedDate.Set(time.Now().UTC())

	newMobileContact1 = models.NewContact(testUser1.ID, "", "555-555-5555", core.Mobile, false)
	newEmailContact1.ConfirmedDate.Set(time.Now().UTC())

	newMobileContact2 = models.NewContact(testUser1.ID, "", "444-444-4444", core.Mobile, false)

	noMatchingUserContact = models.NewContact(testHarness.IDGenerator(false), "", "email@email.email", core.Mobile, false)

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
	t.Run("SwapPrimaryContacts", func(t *testing.T) {
		_testSwapPrimaryContacts(t, *testHarness.ContactRepo)
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
			contactToAdd:   &newEmailContact3,
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.expectedUserID != tc.contactToAdd.UserID {
					t.Errorf("\tuser id expected: got - %s expected - %s", tc.contactToAdd.UserID, tc.expectedUserID)
				}
				if tc.contactToAdd.ID == "" {
					t.Error("\tcontact id is blank")
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.contactID != contact.ID {
					t.Errorf("\tcontact id expected: got - %s expected - %s", contact.ID, tc.contactID)
				}
				if contact.ID == "" {
					t.Error("\tcontact id is blank")
				}
			}
		})
	}
}

func _testGetPrimaryContactByUserID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		userID            string
		contactType       core.ContactType
		expectedContactID string
		expectedErrorCode string
	}
	testCases := []testCase{
		{
			name:              "GIVEN an existing user id with a primary contact EXPECT the primary contact to be returned",
			userID:            testUser1.ID,
			contactType:       core.Email,
			expectedContactID: newEmailContact1.ID,
		},
		{
			name:              "GIVEN an nonexistant user id EXPECT error no user found",
			userID:            nonExistantUserID,
			contactType:       core.Email,
			expectedErrorCode: coreerrors.ErrCodeNoContactFound,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			contact, err := contactRepo.GetPrimaryContactByUserID(context.TODO(), tc.userID, tc.contactType)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if !contact.IsPrimary {
					t.Errorf("\tcontact returned is not primary contact: %v", contact)
				}
				if tc.expectedContactID != contact.ID {
					t.Errorf("\tcontact id was not expected: got - %s expected - %s", contact.ID, tc.expectedContactID)
				}
				if contact.ID == "" {
					t.Error("\tcontact id is blank")
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
				newEmailContact3.ID,
				newMobileContact1.ID,
				newMobileContact2.ID,
			},
		},
		{
			name:               "GIVEN a user id for a user with no contacts EXPECT no contacts to be returned",
			userID:             testUser2.ID,
			expectedContactIds: []string{},
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsFound := len(contacts)
				if numContactsFound != numExpectedContacts {
					t.Errorf("\tnumber of contacts returned not expected amount: got - %d expected - %d", numContactsFound, numExpectedContacts)
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
		contactType        core.ContactType
		expectedContactIds []string
		expectedErrorCode  string
	}
	testCases := []testCase{
		{
			name:        "GIVEN a valid user id and contact type EXPECT contacts of that type returned",
			userID:      testUser1.ID,
			contactType: core.Email,
			expectedContactIds: []string{
				newEmailContact1.ID,
				newEmailContact2.ID,
				newEmailContact3.ID,
			},
		},
		{
			name:        "GIVEN a valid user id and contact type #2 EXPECT contacts of that type returned",
			userID:      testUser1.ID,
			contactType: core.Mobile,
			expectedContactIds: []string{
				newMobileContact1.ID,
				newMobileContact2.ID,
			},
		},
		{
			name:               "GIVEN a valid user id and contact type that has no contacts of that type EXPECT 0 contacts to be returned",
			userID:             testUser2.ID,
			contactType:        core.Email,
			expectedContactIds: []string{},
		},
		{
			name:               "GIVEN a valid user id and an invalid contact type EXPECT 0 contacts to be returned",
			userID:             testUser1.ID,
			contactType:        "not a valid type",
			expectedContactIds: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			numExpectedContacts := len(tc.expectedContactIds)
			contacts, err := contactRepo.GetContactsByUserIDAndType(context.TODO(), tc.userID, tc.contactType)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				numContactsFound := len(contacts)
				if numContactsFound != numExpectedContacts {
					t.Errorf("\tnumber of contacts returned not expected amount: got - %d expected - %d", numContactsFound, numExpectedContacts)
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
						} else if c.Type != tc.contactType {
							t.Errorf("\tcontact type not expected: got - %s expected - %s", c.Type, tc.contactType)
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.contactToUpdate.Principal != tc.newPrincipal {
					t.Errorf("\tupdated contact principal was not expected: got - %s expected - %s", tc.contactToUpdate.Principal, tc.newPrincipal)
				}
				if tc.markConfirmed {
					if !tc.contactToUpdate.ConfirmedDate.HasValue {
						t.Error("\tupdated contact confirmed date does not have a value even thought test case markConfirmed was true")
					}
					if tc.contactToUpdate.ConfirmedDate.Value.Before(preUpdateTime) {
						t.Errorf("\tupdated contact confirmed date not after pre update time when markConfirmed was set: preUpdateTime - %s confirmedDate - %s", preUpdateTime.Format(time.RFC3339), tc.contactToUpdate.ConfirmedDate.Value.Format(time.RFC3339))
					}
				}
				if tc.contactToUpdate.AuditData.ModifiedOnDate.HasValue &&
					tc.contactToUpdate.AuditData.ModifiedOnDate.Value.Before(preUpdateTime) {
					t.Error("\texpected updated contact ModifiedOnDate to be after the preUpdateTime")
				}
				if tc.contactToUpdate.AuditData.ModifiedByID.HasValue &&
					tc.contactToUpdate.AuditData.ModifiedByID.Value != contactRepoCreatedBy {
					t.Errorf("\tupdated contact ModifiedByID not expected: got %s - expected: %s", tc.contactToUpdate.AuditData.ModifiedByID.Value, contactRepoCreatedBy)
				}
			}
		})
	}
}

func _testGetExistingConfirmedContactsCountByPrincipalAndType(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name                     string
		contactPrincipal         string
		contactType              core.ContactType
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
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if numContacts != tc.expectedNumberOfContacts {
					t.Errorf("\tnumber of contacts returned does not match expected: got - %d expected - %d", numContacts, tc.expectedNumberOfContacts)
				}
			}
		})
	}
}

func _testSwapPrimaryContacts(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name                   string
		previousPrimaryContact *models.Contact
		newPrimaryContact      *models.Contact
		expectedErrorCode      string
	}
	testCases := []testCase{
		{
			name:                   "GIVEN a confirmed contact EXPECT success",
			previousPrimaryContact: &newEmailContact1,
			newPrimaryContact:      &newEmailContact3,
		},
		{
			name:                   "GIVEN a confirmed contact EXPECT success (reset original primary contact)",
			previousPrimaryContact: &newEmailContact3,
			newPrimaryContact:      &newEmailContact1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := contactRepo.SwapPrimaryContacts(context.TODO(), tc.previousPrimaryContact, tc.newPrimaryContact, contactRepoCreatedBy)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if tc.previousPrimaryContact.IsPrimary {
					t.Errorf("\texpected previous primary contact isPimary to be false")
				}
				if !tc.newPrimaryContact.IsPrimary {
					t.Errorf("\texpected new primary contact isPrimary to be true")
				}
				ppc, err := contactRepo.GetContactByID(context.TODO(), tc.previousPrimaryContact.ID)
				if err != nil {
					t.Errorf("\tfailed to retreive new primary contact for evaluation %s: %s", err.GetErrorCode(), err.Error())
				}
				if ppc.IsPrimary {
					t.Error("\tprevious primary contact should no longer be marked as primary")
				}
				npc, err := contactRepo.GetContactByID(context.TODO(), tc.newPrimaryContact.ID)
				if err != nil {
					t.Errorf("\tfailed to retreive new primary contact for evaluation %s: %s", err.GetErrorCode(), err.Error())
				}
				if !npc.IsPrimary {
					t.Error("\tnew primary contact should be marked as primary")
				}
			}
		})
	}
}
