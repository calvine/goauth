package service

import (
	"context"
	"time"

	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/utilities"
)

type loginService struct {
	userRepo     repo.UserRepo
	contactRepo  repo.ContactRepo
	auditLogRepo repo.AuditLogRepo
}

func NewLoginService(userRepo repo.UserRepo, contactRepo repo.ContactRepo, auditLogRepo repo.AuditLogRepo) loginService {
	return loginService{
		userRepo:     userRepo,
		contactRepo:  contactRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, error) {
	user, err := ls.userRepo.GetUserByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue {
		if user.LockedOutUntil.Value.After(now) {
			return models.User{}, errors.NewUserLockedOutError(user.ID, true)
		}
	}
	// is contact confirmed?
	contact, err := ls.contactRepo.GetPrimaryContactByUserId(ctx, user.ID)
	if err != nil {
		return models.User{}, err
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		// TODO: return error that primary contact is not confirmed.
	}
	// check password
	saltedString := utilities.InterleaveStrings(password, user.Salt)
	computedHash, err := utilities.SHA512(saltedString)
	if err != nil {
		// TOOD: return custom error for this case...
	}
	if computedHash != user.PasswordHash {
		// TODO: handle failed login
		// TODO: if password check fails increment consecutive failed login attempts and handle logic to set lockout and reset consecutive attempts
	}
	if user.ConsecutiveFailedLoginAttempts > 0 {
		// TODO: reset consecutive failed login attempts
	}
	return user, nil
}

// func (ls loginService) StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, error) {

// }

// func (ls loginService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, error) {

// }

// func (ls loginService) ResetPassword(ctx context.Context, userId string, newPasswordHash string, newSalt string, initiator string) (bool, error) {

// }
