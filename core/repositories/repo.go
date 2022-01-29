package repo

import (
	"context"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
)

type Repo interface {
	GetType() string

	core.NamedComponent
}

type AuditLogRepo interface {
	LogMessage(ctx context.Context, message models.AuditLog) errors.RichError

	Repo
}

type TokenRepo interface {
	// GetToken retreives a token from a store
	GetToken(ctx context.Context, tokenValue string) (models.Token, errors.RichError)
	// PutToken stores a token in a store
	PutToken(ctx context.Context, token models.Token) errors.RichError
	// DeleteToken deletes a token from a store
	DeleteToken(ctx context.Context, tokenValue string) errors.RichError

	Repo
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
	// GetUserAndContactByConfirmedContact gets the user and the primary contact by a confirmed contact principal and contactType
	GetUserAndContactByConfirmedContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError)

	Repo
}

// TODO: add GetPrimaryContacts method?

type ContactRepo interface {
	// GetContactByContactID gets a contact by its id
	GetContactByID(ctx context.Context, id string) (models.Contact, errors.RichError)
	// GetPrimaryContactByUserID gets a users primary contact by user id
	GetPrimaryContactByUserID(ctx context.Context, userID string, contactType string) (models.Contact, errors.RichError)
	// GetContactsByUserID get all of a users contacts by user id
	GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError)
	// GetContactsByUserIDAndType get all contacts belonging to the user based on the userID of the given type
	GetContactsByUserIDAndType(ctx context.Context, userID string, contactType string) ([]models.Contact, errors.RichError)
	// AddContact adds a user contact
	AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError
	// UpdateContact updates a users contact
	UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError
	// GetExistingConfirmedContactsCountByPrincipalAndType gets the count of all contact with the given principal and type which are confrimed.
	GetExistingConfirmedContactsCountByPrincipalAndType(ctx context.Context, contactType, contactPrincipal string) (int64, errors.RichError)
	// SwapPrimaryContacts takes two contact models and sets the isPimary flag to false for previousPrimaryContact and sets the isPrimary flag to true for newPrimaryContact
	// An important note for this function is that the contacts provided to it MUST be of the same type. this logic is contained in the service that calls this function, but is nontheless critical.
	SwapPrimaryContacts(ctx context.Context, previousPrimaryContact, newPrimaryContact *models.Contact, modifiedBy string) errors.RichError

	Repo
}

type ProfileRepo interface {
	// GetProfileByUserID gets a users profile data by user id
	GetProfileByUserID(ctx context.Context, userID string) (models.Profile, errors.RichError)
	// AddProfile adds a users profile data
	AddProfile(ctx context.Context, profile *models.Profile, createdByID string) errors.RichError
	// UpdateUserProfile updates a users profile data
	UpdateUserProfile(ctx context.Context, profile *models.Profile, modifiedByID string) errors.RichError

	Repo
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

	Repo
}

type AppRepo interface {
	GetAppByID(ctx context.Context, id string) (models.App, errors.RichError)
	GetAppsByOwnerID(ctx context.Context, ownerID string) ([]models.App, errors.RichError)
	GetAppByClientID(ctx context.Context, clientID string) (models.App, errors.RichError)
	GetAppAndScopesByClientID(ctx context.Context, clientID string) (models.App, []models.Scope, errors.RichError)
	AddApp(ctx context.Context, app *models.App, createdBy string) errors.RichError
	UpdateApp(ctx context.Context, app *models.App, modifiedBy string) errors.RichError
	DeleteApp(ctx context.Context, app *models.App, deletedBy string) errors.RichError

	GetScopeByID(ctx context.Context, id string) (models.Scope, errors.RichError)
	GetScopesByAppID(ctx context.Context, appID string) ([]models.Scope, errors.RichError)
	AddScope(ctx context.Context, scope *models.Scope, createdBy string) errors.RichError
	UpdateScope(ctx context.Context, scope *models.Scope, modifiedBy string) errors.RichError
	DeleteScope(ctx context.Context, scope *models.Scope, deletedBy string) errors.RichError

	Repo
}

type JWTSigningMaterialRepo interface {
	GetJWTSigningMaterialByKeyID(ctx context.Context, keyID string) (models.JWTSigningMaterial, errors.RichError)

	// TODO: should this not return expired and disabled material or handle filtering in service? I think it should be handled here, but I am going to stew on it...
	// TODO: change to GetValidJWTSigningMaterialByAlgorithmType which will exclude disabled and expired material...
	GetValidJWTSigningMaterialByAlgorithmType(ctx context.Context, algorithmType models.JSMAlgorithmType) ([]models.JWTSigningMaterial, errors.RichError)
	AddJWTSigningMaterial(ctx context.Context, jwtSigningMaterial *models.JWTSigningMaterial, createdBy string) errors.RichError

	Repo
}
