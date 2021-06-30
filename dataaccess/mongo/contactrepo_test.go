package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
)

var (
	newContact1 = models.Contact{
		Principal:     "testuser1@domain.org",
		IsPrimary:     true,
		Type:          core.CONTACT_TYPE_EMAIL,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newContact2 = models.Contact{
		Principal:        "frank@app.space",
		IsPrimary:        false,
		Type:             core.CONTACT_TYPE_EMAIL,
		ConfirmationCode: nullable.NullableString{HasValue: true, Value: "12345"},
		ConfirmedDate:    nullable.NullableTime{},
	}

	newContact3 = models.Contact{
		Principal:     "555-555-5555",
		IsPrimary:     false,
		Type:          core.CONTACT_TYPE_MOBILE,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}
)

func testMongoContactRepo(t *testing.T, userRepo *userRepo) {
	t.Run("AddContact", func(t *testing.T) {
		_testAddContact(t, userRepo)
	})
	t.Run("GetPrimaryContactByUserId", func(t *testing.T) {
		_testGetPrimaryContactByUserId(t, userRepo)
	})
	t.Run("GetContactsByUserId", func(t *testing.T) {
		_testGetContactsByUserId(t, userRepo)
	})
	t.Run("GetContactByConfirmationCode", func(t *testing.T) {
		_testGetContactByConfirmationCode(t, userRepo)
	})
	t.Run("UpdateContact", func(t *testing.T) {
		_testUpdateContact(t, userRepo)
	})
}

func _testAddContact(t *testing.T, userRepo *userRepo) {
	newContact1.UserID = testUser1.ID
	err := userRepo.AddContact(context.TODO(), &newContact1, testUser1.ID)
	if err != nil {
		t.Error("failed to add contact to user", err)
	}

	newContact2.UserID = testUser1.ID
	err = userRepo.AddContact(context.TODO(), &newContact2, testUser1.ID)
	if err != nil {
		t.Error("failed to add contact to user", err)
	}

	newContact3.UserID = testUser1.ID
	err = userRepo.AddContact(context.TODO(), &newContact3, testUser1.ID)
	if err != nil {
		t.Error("failed to add contact to user", err)
	}
}

func _testGetPrimaryContactByUserId(t *testing.T, userRepo *userRepo) {
	userId := testUser1.ID
	contact, err := userRepo.GetPrimaryContactByUserId(context.TODO(), userId)
	if err != nil {
		t.Error("failed to get primary contact for user", userId, err)
	}
	if contact.UserID != testUser1.ID {
		t.Error("expected contact.UserId and testUser1.Id to match", testUser1.ID, contact.UserID)
	}
	if contact.IsPrimary != true {
		t.Error("expected contact.IsPriamay to be true")
	}
	if contact.Principal != newContact1.Principal {
		t.Error("expected contact.Principal to be equal to newContact1.Principal", newContact1.Principal, contact.Principal)
	}

}

func _testGetContactsByUserId(t *testing.T, userRepo *userRepo) {
	userId := testUser1.ID
	contacts, err := userRepo.GetContactsByUserId(context.TODO(), userId)
	if err != nil {
		t.Error("failed to get contacts by user id", userId, err)
	}
	if len(contacts) != 3 {
		t.Error("wrong number of contacts returned", 3, len(contacts))
	}
}

func _testGetContactByConfirmationCode(t *testing.T, userRepo *userRepo) {
	confirmationCode := newContact2.ConfirmationCode.Value
	expectedPrincipal := newContact2.Principal
	expectedIsPrimary := newContact2.IsPrimary
	contact, err := userRepo.GetContactByConfirmationCode(context.TODO(), confirmationCode)
	if err != nil {
		t.Error("failed to get contact by confirmation code", err)
	}
	if contact.Principal != expectedPrincipal {
		t.Error("unexpected value of Principal", expectedPrincipal, contact.Principal)
	}
	if contact.IsPrimary != expectedIsPrimary {
		t.Error("unexpected value of IsPrimary", expectedIsPrimary, contact.IsPrimary)
	}
	if contact.ConfirmationCode.Value != confirmationCode {
		t.Error("expected retreived confirmation code to match passed in confirmation code.", confirmationCode, contact.ConfirmationCode.Value)
	}
}

func _testUpdateContact(t *testing.T, userRepo *userRepo) {
	modifiedByID := "test update contact"
	preUpdateTime := time.Now().UTC()
	newEmail := "a_different_email@mail.org"
	newConfirmedDate := time.Now().UTC()
	newContact3.Principal = newEmail
	newContact3.ConfirmedDate = nullable.NullableTime{}
	newContact3.ConfirmedDate.Set(newConfirmedDate)

	err := userRepo.UpdateContact(context.TODO(), &newContact3, modifiedByID)
	if err != nil {
		t.Error("failed to update contact", err)
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
