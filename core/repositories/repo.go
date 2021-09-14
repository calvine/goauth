package repo

import (
	"context"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
)

//TODO: Make "ID" consistent through code base...
type AuditLogRepo interface {
	LogMessage(ctx context.Context, message models.AuditLog) errors.RichError
}

type TokenRepo interface {
	// GetToken retreives a token from a store
	GetToken(tokenValue string) (models.Token, errors.RichError)
	// PutToken stores a token in a store
	PutToken(token models.Token) errors.RichError
	// DeleteToken deletes a token from a store
	DeleteToken(tokenValue string) errors.RichError
}

// UserRepo is responsible for accessing user data from the database.
type UserRepo interface {
	// GetUserByID gets a user by its id
	GetUserByID(ctx context.Context, id string) (models.User, errors.RichError)
	// AddUser adds a user record
	AddUser(ctx context.Context, user *models.User, createdByID string) errors.RichError
	// UpdateUser updates a user record
	UpdateUser(ctx context.Context, user *models.User, modifiedByID string) errors.RichError
	// GetUserByPrimaryContact gets the user by their primary contact
	GetUserByPrimaryContact(ctx context.Context, contactPrincipalType, contactPrincipal string) (models.User, errors.RichError)
	// GetUserAndContactByPrimaryContact gets the user and the primary contact by their primary contact principal and contactType
	GetUserAndContactByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError)
}

type ContactRepo interface {
	// GetContactByContactID gets a contact by its id
	GetContactByID(ctx context.Context, id string) (models.Contact, errors.RichError)
	// GetPrimaryContactByUserID gets a users primary contact by user id
	GetPrimaryContactByUserID(ctx context.Context, userID string) (models.Contact, errors.RichError)
	// GetContactsByUserID get all of a users contacts by user id
	GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError)
	// GetContactByConfirmationCode get user contact by confirmation code
	// GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, errors.RichError)
	// AddContact adds a user contact
	AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError
	// UpdateContact updates a users contact
	UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError
	// ConfirmContact sets a contact to confirmed based on the received confirmation code.
	// ConfirmContact(ctx context.Context, confirmationCode, modifiedByID string) errors.RichError
}

type ProfileRepo interface {
	// GetProfileByUserID gets a users profile data by user id
	GetProfileByUserID(ctx context.Context, userID string) (models.Profile, errors.RichError)
	// AddProfile adds a users profile data
	AddProfile(ctx context.Context, profile *models.Profile, createdByID string) errors.RichError
	// UpdateUserProfile updates a users profile data
	UpdateUserProfile(ctx context.Context, profile *models.Profile, modifiedByID string) errors.RichError
}

type AddressRepo interface {
	// GetAddressByID gets an address by id
	GetAddressByID(ctx context.Context, id string) (models.Address, errors.RichError)
	// GetPrimaryAddressByUserID gets the prinmary address of a user by the user id
	GetPrimaryAddressByUserID(ctx context.Context, userID string) (models.Address, errors.RichError)
	// GetAddressesByUserID gets address of the user by id
	GetAddressesByUserID(ctx context.Context, userID string) ([]models.Address, errors.RichError)
	// AddAddress adds a user address
	AddAddress(ctx context.Context, address *models.Address, createdByID string) errors.RichError
	// UpdateAddress updates a users address
	UpdateAddress(ctx context.Context, address *models.Address, modifiedByID string) errors.RichError
}

type AppRepo interface {
	GetAppByID(ctx context.Context, id string) (models.App, errors.RichError)
	GetAppsByOwnerID(ctx context.Context, ownerID string) ([]models.App, errors.RichError)
	GetAppAndScopesByClientIDAndCallbackURI(ctx context.Context, clientID, callbackURI string) (models.App, []models.Scope, errors.RichError)
	AddApp(ctx context.Context, app *models.App, createdBy string) errors.RichError
	UpdateApp(ctx context.Context, app *models.App, modifiedBy string) errors.RichError
	DeleteApp(ctx context.Context, app *models.App, deletedBy string) errors.RichError

	GetScopesByAppID(ctx context.Context, appID string) ([]models.Scope, errors.RichError)
	AddScope(ctx context.Context, scope *models.Scope, createdBy string) errors.RichError
	UpdateScope(ctx context.Context, scope *models.Scope, modifiedBy string) errors.RichError
	DeleteScope(ctx context.Context, scope *models.Scope, deletedBy string) errors.RichError
}
