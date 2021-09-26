package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
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

func (contactRepo) GetName() string {
	return "contactRepo"
}

func (contactRepo) GetType() string {
	return dataSourceType
}

func (cr contactRepo) GetContactByID(ctx context.Context, id string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetContactByID", cr.GetType())
	defer span.End()
	contact, ok := (*cr.contacts)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		err := coreerrors.NewNoAppFoundError(fields, true)
		evtString := fmt.Sprintf("contact id not found: %s", id)
		apptelemetry.SetSpanError(&span, err, evtString)
		return contact, err
	}
	span.AddEvent("contact found")
	return contact, nil
}

func (cr contactRepo) GetPrimaryContactByUserID(ctx context.Context, userID string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetPrimaryContactByUserID", cr.GetType())
	defer span.End()
	var contact models.Contact
	contactFound := false
	for _, c := range *cr.contacts {
		if c.UserID == userID && c.IsPrimary {
			contact = c
			contactFound = true
			break
		}
	}
	if !contactFound {
		fields := map[string]interface{}{
			"UserID":    userID,
			"IsPrimary": true,
		}
		err := coreerrors.NewNoContactFoundError(fields, true)
		evtString := fmt.Sprintf("no primary contact found for user: %s", userID)
		apptelemetry.SetSpanError(&span, err, evtString)
		return contact, err
	}
	span.AddEvent("primary contact found")
	return contact, nil
}

func (cr contactRepo) GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetContactsByUserID", cr.GetType())
	defer span.End()
	contacts := make([]models.Contact, 0)
	for _, c := range *cr.contacts {
		if c.UserID == userID {
			contacts = append(contacts, c)
		}
	}
	if len(contacts) == 0 {
		fields := map[string]interface{}{
			"UserID": userID,
		}
		err := coreerrors.NewNoContactFoundError(fields, true)
		evtString := fmt.Sprintf("unable to find contact for user: %s", userID)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return nil, err
	}
	span.AddEvent("contacts found")
	return contacts, nil
}

func (cr contactRepo) AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "AddContact", cr.GetType())
	defer span.End()
	contact.AuditData.CreatedByID = createdByID
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	if contact.ID == "" {
		contact.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*cr.contacts)[contact.ID] = *contact
	span.AddEvent("contact added")
	return nil
}

func (cr contactRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "UpdateContact", cr.GetType())
	defer span.End()
	contact.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedByID}
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*cr.contacts)[contact.ID] = *contact
	span.AddEvent("contact updated")
	return nil
}
