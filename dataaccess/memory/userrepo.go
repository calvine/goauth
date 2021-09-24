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

type userRepo struct {
	users    *map[string]models.User
	contacts *map[string]models.Contact
}

func NewMemoryUserRepo() repo.UserRepo {
	if users == nil {
		users = make(map[string]models.User)
	}
	if contacts == nil {
		contacts = make(map[string]models.Contact)
	}
	return userRepo{
		users:    &users,
		contacts: &contacts,
	}
}

func (userRepo) GetName() string {
	return "userRepo"
}

func (userRepo) GetType() string {
	return dataSourceType
}

func (ur userRepo) GetUserByID(ctx context.Context, id string) (models.User, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ur.GetName()).Start(ctx, "GetUserByID")
	span.SetAttributes(attribute.String("db", ur.GetType()))
	defer span.End()
	user, ok := (*ur.users)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		return user, coreerrors.NewNoUserFoundError(fields, true)
	}
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdByID string) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ur.GetName()).Start(ctx, "AddUser")
	span.SetAttributes(attribute.String("db", ur.GetType()))
	defer span.End()
	user.AuditData.CreatedByID = createdByID
	user.AuditData.CreatedOnDate = time.Now().UTC()
	if user.ID == "" {
		user.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*ur.users)[user.ID] = *user
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedByID string) errors.RichError {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ur.GetName()).Start(ctx, "UpdateUser")
	span.SetAttributes(attribute.String("db", ur.GetType()))
	defer span.End()
	user.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedByID}
	user.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*ur.users)[user.ID] = *user
	return nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactPrincipalType, contactPrincipal string) (models.User, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ur.GetName()).Start(ctx, "GetUserByPrimaryContact")
	span.SetAttributes(attribute.String("db", ur.GetType()))
	defer span.End()
	var user models.User
	var contact models.Contact
	contactFound := false
	for _, c := range *ur.contacts {
		if c.Principal == contactPrincipal &&
			c.Type == contactPrincipalType &&
			c.IsPrimary {
			contact = c
			contactFound = true
			break
		}
	}
	if !contactFound {
		fields := map[string]interface{}{
			"contacts.isPrimary": true,
			"contacts.type":      contactPrincipalType,
			"contacts.principal": contactPrincipal,
		}
		return user, coreerrors.NewNoUserFoundError(fields, true)
	}
	user, ok := (*ur.users)[contact.UserID]
	if !ok {
		// this should not be able to happen...
		fields := map[string]interface{}{
			"contacts.isPrimary": true,
			"contacts.type":      contactPrincipalType,
			"contacts.principal": contactPrincipal,
		}
		return user, coreerrors.NewNoUserFoundError(fields, true)
	}
	return user, nil
}

func (ur userRepo) GetUserAndContactByContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError) {
	spanContext := trace.SpanFromContext(ctx)
	_, span := spanContext.TracerProvider().Tracer(ur.GetName()).Start(ctx, "GetUserAndContactByContact")
	span.SetAttributes(attribute.String("db", ur.GetType()))
	defer span.End()
	var user models.User
	var contact models.Contact
	contactFound := false
	for _, c := range *ur.contacts {
		if c.Principal == contactPrincipal &&
			c.Type == contactType {
			contact = c
			contactFound = true
			break
		}
	}
	if !contactFound {
		fields := map[string]interface{}{
			"contacts.type":      contactType,
			"contacts.principal": contactPrincipal,
		}
		return user, contact, coreerrors.NewNoUserFoundError(fields, true)
	}
	user, ok := (*ur.users)[contact.UserID]
	if !ok {
		// this should not be able to happen...
		fields := map[string]interface{}{
			"contacts.type":      contactType,
			"contacts.principal": contactPrincipal,
		}
		return user, contact, coreerrors.NewNoUserFoundError(fields, true)
	}
	return user, contact, nil
}
