package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	"github.com/calvine/goauth/core/constants/contact"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

type userRepo struct {
	users    *map[string]models.User
	contacts *map[string]models.Contact
}

func NewMemoryUserRepo(users *map[string]models.User, contacts *map[string]models.Contact) (repo.UserRepo, errors.RichError) {
	if users == nil {
		return userRepo{}, coreerrors.NewNilNotAllowedError(true)
	}
	if contacts == nil {
		return userRepo{}, coreerrors.NewNilNotAllowedError(true)
	}
	return userRepo{
		users:    users,
		contacts: contacts,
	}, nil
}

func (userRepo) GetName() string {
	return "userRepo"
}

func (userRepo) GetType() string {
	return dataSourceType
}

func (ur userRepo) GetUserByID(ctx context.Context, id string) (models.User, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserByID", ur.GetType())
	defer span.End()
	user, ok := (*ur.users)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with id: %s", id)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return user, err
	}
	span.AddEvent("retreived user")
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "AddUser", ur.GetType())
	defer span.End()
	user.AuditData.CreatedByID = createdByID
	user.AuditData.CreatedOnDate = time.Now().UTC()
	if user.ID == "" {
		user.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*ur.users)[user.ID] = *user
	span.AddEvent("user added")
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "UpdateUser", ur.GetType())
	defer span.End()
	user.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedByID}
	user.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*ur.users)[user.ID] = *user
	span.AddEvent("user updated")
	return nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactPrincipalType contact.Type, contactPrincipal string) (models.User, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserByPrimaryContact", ur.GetType())
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
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with primary contact %s of type %s", contactPrincipal, contactPrincipalType)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return user, err
	}
	user, ok := (*ur.users)[contact.UserID]
	if !ok {
		// this should not be able to happen...
		fields := map[string]interface{}{
			"contacts.isPrimary": true,
			"contacts.type":      contactPrincipalType,
			"contacts.principal": contactPrincipal,
		}
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with primary contact %s of type %s", contactPrincipal, contactPrincipalType)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return user, err
	}
	span.AddEvent("user and primary contact retreived")
	return user, nil
}

func (ur userRepo) GetUserAndContactByConfirmedContact(ctx context.Context, contactType contact.Type, contactPrincipal string) (models.User, models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserAndContactByConfirmedContact", ur.GetType())
	defer span.End()
	var user models.User
	var contact models.Contact
	contactFound := false
	for _, c := range *ur.contacts {
		if c.Principal == contactPrincipal &&
			c.Type == contactType &&
			c.IsConfirmed() {
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
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with contact %s of type %s", contactPrincipal, contactType)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return user, contact, err
	}
	user, ok := (*ur.users)[contact.UserID]
	if !ok {
		// this should not be able to happen...
		fields := map[string]interface{}{
			"contacts.type":      contactType,
			"contacts.principal": contactPrincipal,
		}
		err := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with contact %s of type %s", contactPrincipal, contactType)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return user, contact, err
	}
	span.AddEvent("user and contact retreived")
	return user, contact, nil
}
