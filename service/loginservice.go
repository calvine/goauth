package service

import (
	"context"
	"time"

	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"golang.org/x/crypto/bcrypt"
)

type loginService struct {
	auditLogRepo repo.AuditLogRepo
	contactRepo  repo.ContactRepo
	emailService EmailService
	userRepo     repo.UserRepo
}

func NewLoginService(auditLogRepo repo.AuditLogRepo, contactRepo repo.ContactRepo, emailService EmailService, userRepo repo.UserRepo) loginService {
	return loginService{
		auditLogRepo: auditLogRepo,
		contactRepo:  contactRepo,
		emailService: emailService,
		userRepo:     userRepo,
	}
}

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, errors.RichError) {
	user, err := ls.userRepo.GetUserByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && user.LockedOutUntil.Value.After(now) {
		return models.User{}, errors.NewUserLockedOutError(user.ID, true)
	}
	// is contact confirmed?
	contact, err := ls.contactRepo.GetPrimaryContactByUserId(ctx, user.ID)
	if err != nil {
		return models.User{}, err
	}
	if !contact.IsPrimary {
		return models.User{}, errors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		// TODO: return error that primary contact is not confirmed.
		return models.User{}, errors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
	}
	// check password
	// saltedString := utilities.InterleaveStrings(password, user.Salt)
	// TODO: use bcrypt...
	// hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if bcryptErr == bcrypt.ErrMismatchedHashAndPassword {
		user.ConsecutiveFailedLoginAttempts += 1
		// TODO: make max ConsecutiveFailedLoginAttempts configurable
		if user.ConsecutiveFailedLoginAttempts >= 10 {
			user.ConsecutiveFailedLoginAttempts = 0
			// TODO: make lockout time configurable
			user.LockedOutUntil.Set(now.Add(time.Minute * 15))
		}
		_ = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		return models.User{}, errors.NewLoginFailedWrongPasswordError(user.ID, true)
	} else if bcryptErr != nil {
		return models.User{}, errors.NewBcryptPasswordHashErrorError(user.ID, bcryptErr, true)
	}
	if user.ConsecutiveFailedLoginAttempts > 0 {
		// reset consecutive failed login attempts because we have a successful login
		user.ConsecutiveFailedLoginAttempts = 0
		err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		if err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}

// func (ls loginService) StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, error) {

// }

// func (ls loginService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, error) {

// }

// func (ls loginService) ResetPassword(ctx context.Context, userId string, newPasswordHash string, newSalt string, initiator string) (bool, error) {

// }
