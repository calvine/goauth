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
	users    *map[string]models.User
	contacts *map[string]models.Contact
}

func NewMemoryContactRepo(users *map[string]models.User, contacts *map[string]models.Contact) (repo.ContactRepo, errors.RichError) {
	if users == nil {
		return contactRepo{}, coreerrors.NewNilParameterNotAllowedError("NewMemoryContactRepo", "users", true)
	}
	if contacts == nil {
		return contactRepo{}, coreerrors.NewNilNotAllowedError(true)
	}
	return contactRepo{
		users:    users,
		contacts: contacts,
	}, nil
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
		err := coreerrors.NewNoContactFoundError(fields, true)
		evtString := fmt.Sprintf("contact id not found: %s", id)
		apptelemetry.SetSpanError(&span, err, evtString)
		return contact, err
	}
	span.AddEvent("contact found")
	return contact, nil
}

func (cr contactRepo) GetPrimaryContactByUserID(ctx context.Context, userID string, contactType string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetPrimaryContactByUserID", cr.GetType())
	defer span.End()
	var contact models.Contact
	contactFound := false
	for _, c := range *cr.contacts {
		if c.UserID == userID &&
			c.IsPrimary &&
			c.Type == contactType {
			contact = c
			contactFound = true
			break
		}
	}
	if !contactFound {
		// what if there is no primary contact of the type found? new error?
		_, ok := (*cr.users)[userID]
		if !ok {
			fields := map[string]interface{}{
				"UserID":    userID,
				"IsPrimary": true,
			}
			err := coreerrors.NewNoContactFoundError(fields, true)
			evtString := fmt.Sprintf("no primary contact found for user: %s", userID)
			apptelemetry.SetSpanError(&span, err, evtString)
			return contact, err
		}
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
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("unable to find contact for user: %s", userID)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return nil, err
	}
	span.AddEvent("contacts found")
	return contacts, nil
}

func (cr contactRepo) GetContactsByUserIDAndType(ctx context.Context, userID string, contactType string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetContactsByUserID", cr.GetType())
	defer span.End()
	contacts := make([]models.Contact, 0)
	for _, c := range *cr.contacts {
		if c.UserID == userID &&
			c.Type == contactType {
			contacts = append(contacts, c)
		}
	}
	if len(contacts) == 0 {
		_, ok := (*cr.users)[userID]
		if !ok {
			// if we get here the user does not exist
			fields := map[string]interface{}{
				"UserID": userID,
			}
			err := coreerrors.NewNoUserFoundError(fields, true)
			evtString := fmt.Sprintf("unable to find contact for user: %s", userID)
			apptelemetry.SetSpanOriginalError(&span, err, evtString)
			return nil, err
		}
	}
	span.AddEvent("contacts retreived")
	return contacts, nil
}

func (cr contactRepo) AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "AddContact", cr.GetType())
	defer span.End()
	_, userFound := (*cr.users)[contact.UserID]
	if !userFound {
		fields := map[string]interface{}{
			"UserID":    contact.UserID,
			"IsPrimary": true,
		}
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found for id: %s", contact.UserID)
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
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

func (cr contactRepo) GetExistingConfirmedContactsCountByPrincipalAndType(ctx context.Context, contactType, contactPrincipal string) (int64, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "GetExistingConfirmedContactsCountByPrincipalAndType", cr.GetType())
	defer span.End()
	numConfirmedContacts := int64(0)
	for _, c := range *cr.contacts {
		if c.IsConfirmed() &&
			c.Type == contactType &&
			c.Principal == contactPrincipal {
			numConfirmedContacts++
		}
	}
	span.AddEvent("number of confirmed contacts retreived")
	return numConfirmedContacts, nil
}

func (cr contactRepo) SwapPrimaryContacts(ctx context.Context, previousPrimaryContact, newPrimaryContact *models.Contact, modifiedBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, cr.GetName(), "SwapPrimaryContacts", cr.GetType())
	defer span.End()
	previousPrimaryContact.IsPrimary = false
	previousPrimaryContact.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	previousPrimaryContact.AuditData.ModifiedByID.Set(modifiedBy)
	(*cr.contacts)[previousPrimaryContact.ID] = *previousPrimaryContact
	newPrimaryContact.IsPrimary = true
	newPrimaryContact.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	newPrimaryContact.AuditData.ModifiedByID.Set(modifiedBy)
	(*cr.contacts)[newPrimaryContact.ID] = *newPrimaryContact
	span.AddEvent("contact primary states set")
	return nil
}
