package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
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
		Principal:        "555-555-5555",
		IsPrimary:        false,
		Type:             core.CONTACT_TYPE_MOBILE,
		ConfirmationCode: nullable.NullableString{HasValue: false},
		ConfirmedDate:    nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	newContact4 = models.Contact{
		Principal:        "email@email.email",
		IsPrimary:        false,
		Type:             core.CONTACT_TYPE_MOBILE,
		ConfirmationCode: nullable.NullableString{HasValue: true, Value: "54321"},
		ConfirmedDate:    nullable.NullableTime{HasValue: false},
	}
)

// func TestInterface(t *testing.T) {
// 	cr := NewUserRepo(nil)
// 	reflect.ValueOf(cr).
// }

func testMongoContactRepo(t *testing.T, contactRepo repo.ContactRepo) {
	t.Run("AddContact", func(t *testing.T) {
		_testAddContact(t, contactRepo)
	})
	t.Run("GetContactById", func(t *testing.T) {
		_testGetContactById(t, contactRepo)
	})
	t.Run("GetPrimaryContactByUserId", func(t *testing.T) {
		_testGetPrimaryContactByUserId(t, contactRepo)
	})
	t.Run("GetContactsByUserId", func(t *testing.T) {
		_testGetContactsByUserId(t, contactRepo)
	})
	t.Run("GetContactByConfirmationCode", func(t *testing.T) {
		_testGetContactByConfirmationCode(t, contactRepo)
	})
	t.Run("UpdateContact", func(t *testing.T) {
		_testUpdateContact(t, contactRepo)
	})
	t.Run("ConfirmContact", func(t *testing.T) {
		_testConfirmContact(t, contactRepo)
	})
}

func _testAddContact(t *testing.T, userRepo repo.ContactRepo) {
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

	newContact4.UserID = testUser1.ID
	err = userRepo.AddContact(context.TODO(), &newContact4, testUser1.ID)
	if err != nil {
		t.Error("failed to add contact to user", err)
	}
}

func _testGetContactById(t *testing.T, userRepo repo.ContactRepo) {
	_, err := userRepo.GetContactByContactId(context.TODO(), newContact1.ID)
	if err != nil {
		t.Error("failed to get contact by given id", newContact1.ID, err)
	}
}

func _testGetPrimaryContactByUserId(t *testing.T, userRepo repo.ContactRepo) {
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

func _testGetContactsByUserId(t *testing.T, userRepo repo.ContactRepo) {
	userId := testUser1.ID
	contacts, err := userRepo.GetContactsByUserId(context.TODO(), userId)
	if err != nil {
		t.Error("failed to get contacts by user id", userId, err)
	}
	if len(contacts) != 4 {
		t.Error("wrong number of contacts returned", 3, len(contacts))
	}
}

func _testGetContactByConfirmationCode(t *testing.T, userRepo repo.ContactRepo) {
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

func _testUpdateContact(t *testing.T, userRepo repo.ContactRepo) {
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

func _testConfirmContact(t *testing.T, contactRepo repo.ContactRepo) {
	modifiedById := "test confirm contact"
	preConfirmTime := time.Now().UTC()
	time.Sleep(time.Second * 3)
	err := contactRepo.ConfirmContact(context.TODO(), newContact4.ConfirmationCode.Value, modifiedById)
	if err != nil {
		t.Error("failed to update contact", err)
	}
	updatedContact4, err := contactRepo.GetContactByContactId(context.TODO(), newContact4.ID)
	if err != nil {
		t.Error("failed to get instance of updated contact for comparison", err)
	}
	if updatedContact4.ConfirmationCode.HasValue {
		t.Error("confirmed contact confirmation code should not have value", updatedContact4.ConfirmationCode.Value)
	}
	if !updatedContact4.ConfirmedDate.HasValue {
		t.Error("confirmed contact confirmed date should have a value")
	} else if updatedContact4.ConfirmedDate.Value.Before(preConfirmTime) {
		t.Error("confirmed contact confirmed date should be after the pre confirm time", preConfirmTime.String(), updatedContact4.ConfirmedDate.Value.String())
	}
	if updatedContact4.AuditData.ModifiedByID.Value != modifiedById {
		t.Error("confirmed contact modified by id should match modifiedById", modifiedById, updatedContact4.AuditData.ModifiedByID.Value)
	}
	if updatedContact4.AuditData.ModifiedOnDate.Value.Before(preConfirmTime) {
		t.Error("confirmed contact modified on date should be after the pre confirm time", preConfirmTime.String(), updatedContact4.AuditData.ModifiedOnDate.Value.String())
	}
}
