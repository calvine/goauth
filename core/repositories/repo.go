package repo

import (
	"context"

	"github.com/calvine/goauth/core/models"
)

// TODO: change all instances of error to RichError

type AuditLogRepo interface {
	LogMessage(ctx context.Context, message models.AuditLog) error
}

// UserRepo is responsible for accessing user data from the database.
type UserRepo interface {
	// GetUserById gets a user by its id
	GetUserById(ctx context.Context, id string) (models.User, error)
	// AddUser adds a user record
	AddUser(ctx context.Context, user *models.User, createdById string) error
	// UpdateUser updates a user record
	UpdateUser(ctx context.Context, user *models.User, modifiedById string) error
	// GetUserByPrimaryContact gets the user by their primary contact
	GetUserByPrimaryContact(ctx context.Context, contactPrincipalType, contactPrincipal string) (models.User, error)
}

type ContactRepo interface {
	// GetContactByContactId gets a contact by its id
	GetContactByContactId(ctx context.Context, contactId string) (models.Contact, error)
	// GetPrimaryContactByUserId gets a users primary contact by user id
	GetPrimaryContactByUserId(ctx context.Context, userId string) (models.Contact, error)
	// GetContactsByUserId get all of a users contacts by user id
	GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, error)
	// GetContactByConfirmationCode get user contact by confirmation code
	GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, error)
	// AddContact adds a user contact
	AddContact(ctx context.Context, contact *models.Contact, createdById string) error
	// UpdateContact updates a users contact
	UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) error
	// ConfirmContact sets a contact to confirmed based on the received confirmation code.
	ConfirmContact(ctx context.Context, confirmationCode, modifiedById string) error
}

type ProfileRepo interface {
	// GetProfileByUserId gets a users profile data by user id
	GetProfileByUserId(ctx context.Context, userId string) (models.Profile, error)
	// AddProfile adds a users profile data
	AddProfile(ctx context.Context, profile *models.Profile, createdById string) error
	// UpdateUserProfile updates a users profile data
	UpdateUserProfile(ctx context.Context, profile *models.Profile, modifiedById string) error
}

type AddressRepo interface {
	// GetAddressById gets an address by id
	GetAddressById(ctx context.Context, id string) (models.Address, error)
	// GetPrimaryAddressByUserId gets the prinmary address of a user by the user id
	GetPrimaryAddressByUserId(ctx context.Context, userId string) (models.Address, error)
	// GetAddressesByUserId gets address of the user by Id
	GetAddressesByUserId(ctx context.Context, userId string) ([]models.Address, error)
	// AddAddress adds a user address
	AddAddress(ctx context.Context, address *models.Address, createdById string) error
	// UpdateAddress updates a users address
	UpdateAddress(ctx context.Context, address *models.Address, modifiedById string) error
}
