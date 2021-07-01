package services

import (
	"context"

	"github.com/calvine/goauth/core/models"
)

// TODO: implemet account lockout after X number of consecutive failed login attempts.
// LoginService is a service used to facilitate logging in
type LoginService interface {
	// LoginWithContact attempts to confirm a users credentials and if they match it returns true and resets the users ConsecutiveFailedLoginAttempts, otherwise it returns false and increments the users ConsecutiveFailedLoginAttempts
	// The principal should only work when it has been confirmed
	LoginWithContact(ctx context.Context, principal, principalType, password string, initiator string) (bool, error)
	// StartPasswordResetByContact sets a password reset token for the user with the corresponding principal and type that are confirmed.
	StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, error)
	// ConfirmContact takes a confirmation code and updates the users contact record to be confirmed.
	ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, error)
	// ResetPassword resets a users password given a userId and new password hash and salt.
	ResetPassword(ctx context.Context, userId string, newPasswordHash string, newSalt string, initiator string) (bool, error)
}

// UserDataService is a service that facilitates access to user related data.
type UserDataService interface {
	// GetUserByConfirmedContact gets a user record via a confirmed contact
	GetUserByConfirmedContact(ctx context.Context, contactPrincipal string, initiator string) (models.User, error)
	// AddUser adds a user record to the database
	AddUser(ctx context.Context, user *models.User, initiator string) error
	// UpdateUser updates a use record in the database
	UpdateUser(ctx context.Context, user *models.User, initiator string) error
	// GetUserPrimaryContact gets a users primary contact
	GetUserPrimaryContact(ctx context.Context, userId string, initiator string) (models.Contact, error)
	// GetUsersContacts gets all of a users contacts
	GetUsersContacts(ctx context.Context, userId string, initiator string) ([]models.Contact, error)
	// GetUsersConfirmedContacts gets all of a users confirmed contacts
	GetUsersConfirmedContacts(ctx context.Context, userId string, initiator string) ([]models.Contact, error)
	// AddContact adds a contact to a user
	AddContact(ctx context.Context, contact *models.Contact, initiator string) error
	// UpdateContact updates a contact for a user
	UpdateContact(ctx context.Context, contact *models.Contact, initiator string) error
}
