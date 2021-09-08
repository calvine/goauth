package repo

import (
	"context"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
)

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
	// GetUserById gets a user by its id
	GetUserById(ctx context.Context, id string) (models.User, errors.RichError)
	// AddUser adds a user record
	AddUser(ctx context.Context, user *models.User, createdById string) errors.RichError
	// UpdateUser updates a user record
	UpdateUser(ctx context.Context, user *models.User, modifiedById string) errors.RichError
	// GetUserByPrimaryContact gets the user by their primary contact
	GetUserByPrimaryContact(ctx context.Context, contactPrincipalType, contactPrincipal string) (models.User, errors.RichError)
	// GetUserAndContactByPrimaryContact gets the user and the primary contact by their primary contact principal and contactType
	GetUserAndContactByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError)
}

type ContactRepo interface {
	// GetContactByContactId gets a contact by its id
	GetContactByContactId(ctx context.Context, contactId string) (models.Contact, errors.RichError)
	// GetPrimaryContactByUserId gets a users primary contact by user id
	GetPrimaryContactByUserId(ctx context.Context, userId string) (models.Contact, errors.RichError)
	// GetContactsByUserId get all of a users contacts by user id
	GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, errors.RichError)
	// GetContactByConfirmationCode get user contact by confirmation code
	// GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, errors.RichError)
	// AddContact adds a user contact
	AddContact(ctx context.Context, contact *models.Contact, createdById string) errors.RichError
	// UpdateContact updates a users contact
	UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) errors.RichError
	// ConfirmContact sets a contact to confirmed based on the received confirmation code.
	// ConfirmContact(ctx context.Context, confirmationCode, modifiedById string) errors.RichError
}

type ProfileRepo interface {
	// GetProfileByUserId gets a users profile data by user id
	GetProfileByUserId(ctx context.Context, userId string) (models.Profile, errors.RichError)
	// AddProfile adds a users profile data
	AddProfile(ctx context.Context, profile *models.Profile, createdById string) errors.RichError
	// UpdateUserProfile updates a users profile data
	UpdateUserProfile(ctx context.Context, profile *models.Profile, modifiedById string) errors.RichError
}

type AddressRepo interface {
	// GetAddressById gets an address by id
	GetAddressById(ctx context.Context, id string) (models.Address, errors.RichError)
	// GetPrimaryAddressByUserId gets the prinmary address of a user by the user id
	GetPrimaryAddressByUserId(ctx context.Context, userId string) (models.Address, errors.RichError)
	// GetAddressesByUserId gets address of the user by Id
	GetAddressesByUserId(ctx context.Context, userId string) ([]models.Address, errors.RichError)
	// AddAddress adds a user address
	AddAddress(ctx context.Context, address *models.Address, createdById string) errors.RichError
	// UpdateAddress updates a users address
	UpdateAddress(ctx context.Context, address *models.Address, modifiedById string) errors.RichError
}
