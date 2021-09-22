package services

import (
	"context"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

// TODO: change all instances of error to RichError

// LoginService is a service used to facilitate logging in
type LoginService interface {
	// LoginWithContact attempts to confirm a users credentials and if they match it returns true and resets the users ConsecutiveFailedLoginAttempts, otherwise it returns false and increments the users ConsecutiveFailedLoginAttempts
	// The principal should only work when it has been confirmed
	LoginWithPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType, password string, initiator string) (models.User, errors.RichError)
	// StartPasswordResetByContact sets a password reset token for the user with the corresponding principal and type that are confirmed.
	StartPasswordResetByPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType string, initiator string) (string, errors.RichError)
	// ResetPassword resets a users password given a password reset token and new password hash and salt.
	ResetPassword(ctx context.Context, logger *zap.Logger, passwordResetToken string, newPassword string, initiator string) errors.RichError
}

// UserService is a service that facilitates access to user related data.
type UserService interface {
	// GetUserByConfirmedContact gets a user record via a confirmed contact
	GetUserByConfirmedContact(ctx context.Context, contactPrincipal string, initiator string) (models.User, errors.RichError)
	// AddUser adds a user record to the database
	AddUser(ctx context.Context, user *models.User, initiator string) errors.RichError
	// UpdateUser updates a use record in the database
	UpdateUser(ctx context.Context, user *models.User, initiator string) errors.RichError
	// GetUserPrimaryContact gets a users primary contact
	GetUserPrimaryContact(ctx context.Context, userID string, initiator string) (models.Contact, errors.RichError)
	// GetUsersContacts gets all of a users contacts
	GetUsersContacts(ctx context.Context, userID string, initiator string) ([]models.Contact, errors.RichError)
	// GetUsersConfirmedContacts gets all of a users confirmed contacts
	GetUsersConfirmedContacts(ctx context.Context, userID string, initiator string) ([]models.Contact, errors.RichError)
	// AddContact adds a contact to a user
	AddContact(ctx context.Context, contact *models.Contact, initiator string) errors.RichError
	// UpdateContact updates a contact for a user
	UpdateContact(ctx context.Context, contact *models.Contact, initiator string) errors.RichError
	// ConfirmContact takes a confirmation code and updates the users contact record to be confirmed.
	ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, errors.RichError)
}

type AppService interface {
	GetAppsByOwnerID(ctx context.Context, logger *zap.Logger, ownerID string, initiator string) ([]models.App, errors.RichError)
	GetAppByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.App, errors.RichError)
	GetAppByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, errors.RichError)
	GetAppAndScopesByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, []models.Scope, errors.RichError)
	AddApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	UpdateApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	DeleteApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	GetScopeByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.Scope, errors.RichError)
	GetScopesByAppID(ctx context.Context, logger *zap.Logger, appID string, initiator string) ([]models.Scope, errors.RichError)
	// TODO: Determine if needed...
	// GetScopesByClientID(ctx context.Context, clientID string, initiator string) ([]models.Scope, errors.RichError)

	AddScopeToApp(ctx context.Context, logger *zap.Logger, scopes *models.Scope, initiator string) errors.RichError
	UpdateScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError
	DeleteScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError
}

type EmailService interface {
	SendPlainTextEmail(ctx context.Context, to []string, subject, body string) errors.RichError
}

type TokenService interface {
	GetToken(ctx context.Context, tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError)
	PutToken(ctx context.Context, token models.Token) errors.RichError
	DeleteToken(ctx context.Context, tokenValue string) errors.RichError
}
