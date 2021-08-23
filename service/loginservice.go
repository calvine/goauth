package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
	"golang.org/x/crypto/bcrypt"
)

type loginService struct {
	auditLogRepo repo.AuditLogRepo
	contactRepo  repo.ContactRepo
	emailService EmailService
	userRepo     repo.UserRepo
}

func NewLoginService(auditLogRepo repo.AuditLogRepo, contactRepo repo.ContactRepo, emailService EmailService, userRepo repo.UserRepo) services.LoginService {
	return loginService{
		auditLogRepo: auditLogRepo,
		contactRepo:  contactRepo,
		emailService: emailService,
		userRepo:     userRepo,
	}
}

//TODO: Add audit logging

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, errors.RichError) {
	user, contact, err := ls.userRepo.GetUserAndContactByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	if !contact.IsPrimary {
		return models.User{}, coreerrors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && user.LockedOutUntil.Value.After(now) {
		return models.User{}, coreerrors.NewUserLockedOutError(user.ID, true)
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		return models.User{}, coreerrors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
	}
	// check password
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
		return models.User{}, coreerrors.NewLoginFailedWrongPasswordError(user.ID, true)
	} else if bcryptErr != nil {
		return models.User{}, coreerrors.NewBcryptPasswordHashErrorError(user.ID, bcryptErr, true)
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

func (ls loginService) StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, errors.RichError) {
	now := time.Now().UTC()
	user, contact, err := ls.userRepo.GetUserAndContactByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return "", err
	}
	passwordResetToken, roErr := utilities.NewPasswordResetToken()
	if roErr != nil {
		return "", roErr
	}
	user.PasswordResetToken.Set(passwordResetToken)
	// TODO: make password reset token expiration configurable.
	user.PasswordResetTokenExpiration.Set(now.Add(time.Hour * 1))
	err = ls.userRepo.UpdateUser(ctx, &user, initiator)
	if err != nil {
		return "", err
	}
	switch contact.Type {
	case core.CONTACT_TYPE_EMAIL:
		// TODO: create template for this...
		body := fmt.Sprintf("A Password reset has been initiated. Your password reset token is: %s", passwordResetToken)
		ls.emailService.SendPlainTextEmail([]string{contact.Principal}, "Password reset", body)
	default:
		return "", coreerrors.NewComponentNotImplementedError("notification system", fmt.Sprintf("%s notification service", contact.Type), true)
	}
	return passwordResetToken, nil
}

func (ls loginService) ResetPassword(ctx context.Context, passwordResetToken string, newPasswordHash string, initiator string) (bool, errors.RichError) {
	if newPasswordHash == "" {
		return false, coreerrors.NewNoNewPasswordHashProvidedError(true)
	}
	user, err := ls.userRepo.GetUserByPasswordResetToken(ctx, passwordResetToken)
	if err != nil {
		return false, err
	}
	now := time.Now().UTC()
	if now.After(user.PasswordResetTokenExpiration.Value) {
		return false, coreerrors.NewExpiredPasswordExpirationTokenError(passwordResetToken, user.PasswordResetTokenExpiration.Value, true)
	}
	err = ls.userRepo.UpdateUser(ctx, &user, initiator)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ls loginService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, errors.RichError) {
	// TODO: implement this
	return false, nil
}
