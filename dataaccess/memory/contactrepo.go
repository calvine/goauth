package memory

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

type contactRepo struct {
	contacts *map[string]models.Contact
}

func NewMemoryContactRepo() repo.ContactRepo {
	if contacts == nil {
		contacts = make(map[string]models.Contact)
	}
	return contactRepo{
		contacts: &contacts,
	}
}

func (cr contactRepo) GetContactById(ctx context.Context, id string) (models.Contact, errors.RichError) {
	contact, ok := (*cr.contacts)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		return contact, coreerrors.NewNoAppFoundError(fields, true)
	}
	return contact, nil
}

func (cr contactRepo) GetPrimaryContactByUserId(ctx context.Context, userId string) (models.Contact, errors.RichError) {
	var contact models.Contact
	contactFound := false
	for _, c := range *cr.contacts {
		if c.UserID == userId && c.IsPrimary {
			contact = c
			contactFound = true
			break
		}
	}
	if !contactFound {
		fields := map[string]interface{}{
			"UserID":    userId,
			"IsPrimary": true,
		}
		return contact, coreerrors.NewNoContactFoundError(fields, true)
	}
	return contact, nil
}

func (cr contactRepo) GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, errors.RichError) {
	contacts := make([]models.Contact, 0)
	for _, c := range *cr.contacts {
		if c.UserID == userId {
			contacts = append(contacts, c)
		}
	}
	if len(contacts) == 0 {
		fields := map[string]interface{}{
			"UserID": userId,
		}
		return nil, coreerrors.NewNoContactFoundError(fields, true)
	}
	return contacts, nil
}

func (cr contactRepo) AddContact(ctx context.Context, contact *models.Contact, createdById string) errors.RichError {
	contact.AuditData.CreatedByID = createdById
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	(*cr.contacts)[contact.ID] = *contact
	return nil
}

func (cr contactRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) errors.RichError {
	contact.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedById}
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*cr.contacts)[contact.ID] = *contact
	return nil
}
