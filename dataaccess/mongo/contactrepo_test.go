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
	newContact1.UserId = testUser1.Id
	err := userRepo.AddContact(context.TODO(), &newContact1, testUser1.Id)
	if err != nil {
		t.Error("failed to add contact to user", err)
	}
}

func _testGetPrimaryContactByUserId(t *testing.T, userRepo *userRepo) {
	userId := testUser1.Id
	contact, err := userRepo.GetPrimaryContactByUserId(context.TODO(), userId)
	if err != nil {
		t.Error("failed to get primary contact for user", userId, err)
	}
	if contact.UserId != testUser1.Id {
		t.Error("expected contact.UserId and testUser1.Id to match", testUser1.Id, contact.UserId)
	}
	if contact.IsPrimary != true {
		t.Error("expected contact.IsPriamay to be true")
	}
	if contact.Principal != newContact1.Principal {
		t.Error("expected contact.Principal to be equal to newContact1.Principal", newContact1.Principal, contact.Principal)
	}

}

func _testGetContactsByUserId(t *testing.T, userRepo *userRepo) {
	t.Error("not implemented")
}

func _testGetContactByConfirmationCode(t *testing.T, userRepo *userRepo) {
	t.Error("not implemented")
}

func _testUpdateContact(t *testing.T, userRepo *userRepo) {
	t.Error("not implemented")
}
