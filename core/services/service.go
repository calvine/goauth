package services

import (
	"context"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

type Service interface {
	core.NamedComponent
}

// LoginService is a service used to facilitate logging in
type LoginService interface {
	// LoginWithContact attempts to confirm a users credentials and if they match it returns true and resets the users ConsecutiveFailedLoginAttempts, otherwise it returns false and increments the users ConsecutiveFailedLoginAttempts
	// The principal should only work when it has been confirmed
	LoginWithPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType, password string, initiator string) (models.User, errors.RichError)
	// StartPasswordResetByContact sets a password reset token for the user with the corresponding principal and type that are confirmed.
	StartPasswordResetByPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType string, initiator string) (string, errors.RichError)
	// ResetPassword resets a users password given a password reset token and new password hash and salt.
	ResetPassword(ctx context.Context, logger *zap.Logger, passwordResetToken string, newPassword string, initiator string) errors.RichError

	Service
}

// UserService is a service that facilitates access to user related data.
type UserService interface {
	// GetUserAndContactByConfirmedContact gets a user and specified contact record via a confirmed contact
	GetUserAndContactByConfirmedContact(ctx context.Context, logger *zap.Logger, contactType string, contactPrincipal string, initiator string) (models.User, models.Contact, errors.RichError)
	// // AddUser adds a user record to the database
	// AddUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError
	// // UpdateUser updates a use record in the database
	// UpdateUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError
	// RegisterUserAndPrimaryContact registers a new user. it has several responsibilities.
	//	1. ensure no other user has the contact provided as a confirmed contact.
	//	2. send notification to user with link to confirm contact and set password
	RegisterUserAndPrimaryContact(ctx context.Context, logger *zap.Logger, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError)
	// GetUserPrimaryContact gets a users primary contact
	GetUserPrimaryContact(ctx context.Context, logger *zap.Logger, userID string, initiator string) (models.Contact, errors.RichError)
	// GetUsersContacts gets all of a users contacts
	GetUsersContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError)
	// GetUsersConfirmedContacts gets all of a users confirmed contacts
	GetUsersConfirmedContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError)
	// AddContact adds a contact to a user
	AddContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError
	// UpdateContact updates a contact for a user
	UpdateContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError
	// ConfirmContact takes a confirmation code and updates the users contact record to be confirmed.
	ConfirmContact(ctx context.Context, logger *zap.Logger, confirmationCode string, initiator string) errors.RichError

	Service
}

type AppService interface {
	// GetAppsByOwnerID retreives apps beloging to an owner by their id
	GetAppsByOwnerID(ctx context.Context, logger *zap.Logger, ownerID string, initiator string) ([]models.App, errors.RichError)
	// GetAppByID retreives an app by its id
	GetAppByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.App, errors.RichError)
	// GetAppByClientID retreives an app by its client id
	GetAppByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, errors.RichError)
	// GetAppAndScopesByClientID retreives an app and its scoipes by its client id
	GetAppAndScopesByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, []models.Scope, errors.RichError)
	// AddApp adds an app to the underlying data store
	AddApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	// UpdateApp updates an app in the underlying data store
	UpdateApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	// DeleteApp deletes an app
	DeleteApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError
	// GetScopeByID retrives a scope by its id
	GetScopeByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.Scope, errors.RichError)
	// GetScopesByAppID get scopes for an app by its id
	GetScopesByAppID(ctx context.Context, logger *zap.Logger, appID string, initiator string) ([]models.Scope, errors.RichError)
	// TODO: Determine if needed...
	// GetScopesByClientID(ctx context.Context, clientID string, initiator string) ([]models.Scope, errors.RichError)

	// AddScopeToApp adds a scope to the underlying data store for an app given by the scopes app id
	AddScopeToApp(ctx context.Context, logger *zap.Logger, scopes *models.Scope, initiator string) errors.RichError
	// UpdateScope updates a scope in the underlying data store
	UpdateScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError
	// DeleteScope deletes a scope from the underlying data store
	DeleteScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError

	Service
}

type EmailService interface {
	// SendPlainTextEmail sends a plain text email.
	SendPlainTextEmail(ctx context.Context, logger *zap.Logger, to []string, subject, body string) errors.RichError

	Service
}

type TokenService interface {
	// GetToken retreives a token from the underlying data store given its token type and value
	GetToken(ctx context.Context, logger *zap.Logger, tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError)
	// PutToken stores a token in the underlying data store
	PutToken(ctx context.Context, logger *zap.Logger, token models.Token) errors.RichError
	// DeleteToken deletes a token from the underlying data store
	DeleteToken(ctx context.Context, logger *zap.Logger, tokenValue string) errors.RichError

	Service
}
