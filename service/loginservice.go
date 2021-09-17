package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	coreservices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultAccountLockoutDuration time.Duration = time.Minute * 15
	defaultMaxFailedLoginAttempts int           = 10
)

type loginService struct {
	auditLogRepo           repo.AuditLogRepo
	contactRepo            repo.ContactRepo
	emailService           coreservices.EmailService
	userRepo               repo.UserRepo
	tokenService           coreservices.TokenService
	maxFailedLoginAttempts int
	accountLockoutDuration time.Duration
}

type LoginServiceOptions struct {
	AuditLogRepo           repo.AuditLogRepo
	ContactRepo            repo.ContactRepo
	EmailService           coreservices.EmailService
	UserRepo               repo.UserRepo
	TokenService           coreservices.TokenService
	MaxFailedLoginAttempts int
	AccountLockoutDuration time.Duration
}

func NewLoginService(options LoginServiceOptions) coreservices.LoginService {
	if options.MaxFailedLoginAttempts <= 0 {
		options.MaxFailedLoginAttempts = defaultMaxFailedLoginAttempts
	}
	if options.AccountLockoutDuration <= 0 {
		options.AccountLockoutDuration = defaultAccountLockoutDuration
	}
	return loginService{
		auditLogRepo:           options.AuditLogRepo,
		contactRepo:            options.ContactRepo,
		emailService:           options.EmailService,
		userRepo:               options.UserRepo,
		tokenService:           options.TokenService,
		maxFailedLoginAttempts: options.MaxFailedLoginAttempts,
		accountLockoutDuration: options.AccountLockoutDuration,
	}
}

//TODO: Add audit logging

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, errors.RichError) {
	user, contact, err := ls.userRepo.GetUserAndContactByContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && now.Before(user.LockedOutUntil.Value) {
		return models.User{}, coreerrors.NewUserLockedOutError(user.ID, true)
	}
	if !contact.IsPrimary {
		return models.User{}, coreerrors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		return models.User{}, coreerrors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
	}
	// check password
	// hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	passwordMatch, err := utilities.BcryptCompareStringAndHash(user.PasswordHash, password, user.ID)
	if err != nil {
		return models.User{}, err
	}
	if !passwordMatch {
		user.ConsecutiveFailedLoginAttempts += 1

		if user.ConsecutiveFailedLoginAttempts >= ls.maxFailedLoginAttempts {
			user.ConsecutiveFailedLoginAttempts = 0
			user.LockedOutUntil.Set(now.Add(ls.accountLockoutDuration))
		}
		err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		if err != nil {
			// TODO: log this error better!
			fmt.Printf("failed to lock out user %s: %s", user.ID, err.Error())
		}
		return models.User{}, coreerrors.NewLoginFailedWrongPasswordError(user.ID, true)
	}
	if user.ConsecutiveFailedLoginAttempts > 0 {
		// reset consecutive failed login attempts because we have a successful login
		user.ConsecutiveFailedLoginAttempts = 0
		err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		if err != nil {
			return models.User{}, err
		}
	}
	user.LastLoginDate.Set(now)
	user.ConsecutiveFailedLoginAttempts = 0
	user.LockedOutUntil.Unset()
	err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (ls loginService) StartPasswordResetByPrimaryContact(ctx context.Context, principal, principalType string, initiator string) (string, errors.RichError) {
	user, contact, err := ls.userRepo.GetUserAndContactByContact(ctx, principalType, principal)
	if err != nil {
		return "", err
	}
	if !contact.IsPrimary {
		return "", coreerrors.NewPasswordResetContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	// TODO: make password reset token expiration configurable.
	token, err := models.NewToken(user.ID, models.TokenTypePasswordReset, time.Minute*15)
	if err != nil {
		return "", err
	}
	err = ls.tokenService.PutToken(ctx, token)
	if err != nil {
		return "", err
	}
	switch contact.Type {
	case core.CONTACT_TYPE_EMAIL:
		// TODO: create template for this...
		body := fmt.Sprintf("A Password reset has been initiated. Your password reset token is: %s", token.Value)
		ls.emailService.SendPlainTextEmail([]string{contact.Principal}, "Password reset", body)
	default:
		return "", coreerrors.NewComponentNotImplementedError("notification system", fmt.Sprintf("%s notification service", contact.Type), true)
	}
	return token.Value, nil
}

func (ls loginService) ResetPassword(ctx context.Context, passwordResetToken string, newPassword string, initiator string) errors.RichError {
	if newPassword == "" {
		return coreerrors.NewNoNewPasswordHashProvidedError(true)
	}
	// TODO: add password validation logic
	token, err := ls.tokenService.GetToken(ctx, passwordResetToken, models.TokenTypePasswordReset)
	if err != nil {
		return err
	}
	user, err := ls.userRepo.GetUserByID(ctx, token.TargetID)
	if err != nil {
		return err
	}
	newPasswordHash, err := utilities.BcryptHashString(newPassword, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = newPasswordHash
	err = ls.userRepo.UpdateUser(ctx, &user, initiator)
	if err != nil {
		return err
	}
	return nil
}
