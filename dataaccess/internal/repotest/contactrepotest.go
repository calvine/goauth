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
	_, err := contactRepo.GetContactByID(context.TODO(), newContact1.ID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to get contact by given id", newContact1.ID, err.GetErrorCode())
	}
}

func _testGetPrimaryContactByUserID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
	userID := testUser1.ID
	contact, err := contactRepo.GetPrimaryContactByUserID(context.TODO(), userID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to get primary contact for user", userID, err.GetErrorCode())
	}
	if contact.UserID != testUser1.ID {
		t.Error("expected contact.UserID and testUser1.ID to match", testUser1.ID, contact.UserID)
	}
	if !contact.IsPrimary {
		t.Error("expected contact.IsPriamay to be true")
	}
	if contact.Principal != newContact1.Principal {
		t.Error("expected contact.Principal to be equal to newContact1.Principal", newContact1.Principal, contact.Principal)
	}

}

func _testGetContactsByUserID(t *testing.T, contactRepo repo.ContactRepo) {
	type testCase struct {
		name              string
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
	userID := testUser1.ID
	contacts, err := contactRepo.GetContactsByUserID(context.TODO(), userID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to get contacts by user id", userID, err.GetErrorCode())
	}
	if len(contacts) != 4 {
		t.Error("wrong number of contacts returned", 3, len(contacts))
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
		expectedErrorCode string
	}
	testCases := []testCase{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

		})
	}
	modifiedByID := "test update contact"
	preUpdateTime := time.Now().UTC()
	newEmail := "a_different_email@mail.org"
	newConfirmedDate := time.Now().UTC()
	newContact3.Principal = newEmail
	newContact3.ConfirmedDate = nullable.NullableTime{}
	newContact3.ConfirmedDate.Set(newConfirmedDate)

	err := contactRepo.UpdateContact(context.TODO(), &newContact3, modifiedByID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to update contact", err.GetErrorCode())
	}
	if newContact3.Principal != newEmail {
		t.Error("expected principal to be updated", newEmail, newContact3.Principal)
	}
	if !newContact3.ConfirmedDate.HasValue {
		t.Error("expected ConfirmedDate to be have a value")
	}
	if !newContact3.AuditData.ModifiedOnDate.Value.After(preUpdateTime) {
		t.Error("expected ModifiedOnDate to be after the preUpdateTime")
	}
	if newContact3.AuditData.ModifiedByID.Value != modifiedByID {
		t.Error("expected ModifiedByID to be updated", modifiedByID, newContact3.AuditData.ModifiedByID.Value)
	}
}
