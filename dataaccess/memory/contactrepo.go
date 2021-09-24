package memory

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(cr.GetName()).Start(ctx, "GetContactByID")
	span.SetAttributes(attribute.String("db", cr.GetType()))
	defer span.End()
	contact, ok := (*cr.contacts)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		return contact, coreerrors.NewNoAppFoundError(fields, true)
	}
	return contact, nil
}

func (cr contactRepo) GetPrimaryContactByUserID(ctx context.Context, userID string) (models.Contact, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(cr.GetName()).Start(ctx, "GetPrimaryContactByUserID")
	span.SetAttributes(attribute.String("db", cr.GetType()))
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
		return contact, coreerrors.NewNoContactFoundError(fields, true)
	}
	return contact, nil
}

func (cr contactRepo) GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(cr.GetName()).Start(ctx, "GetContactsByUserID")
	span.SetAttributes(attribute.String("db", cr.GetType()))
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
		return nil, coreerrors.NewNoContactFoundError(fields, true)
	}
	return contacts, nil
}

func (cr contactRepo) AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(cr.GetName()).Start(ctx, "AddContact")
	span.SetAttributes(attribute.String("db", cr.GetType()))
	defer span.End()
	contact.AuditData.CreatedByID = createdByID
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	if contact.ID == "" {
		contact.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*cr.contacts)[contact.ID] = *contact
	return nil
}

func (cr contactRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(cr.GetName()).Start(ctx, "UpdateContact")
	span.SetAttributes(attribute.String("db", cr.GetType()))
	defer span.End()
	contact.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedByID}
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*cr.contacts)[contact.ID] = *contact
	return nil
}
