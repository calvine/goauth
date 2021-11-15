package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	coreservices "github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
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

func (loginService) GetName() string {
	return "loginService"
}

// TODO: Add audit logging

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType, password string, initiator string) (models.User, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, ls.GetName(), "LoginWithPrimaryContact")
	defer span.End()
	user, contact, err := ls.userRepo.GetUserAndContactByConfirmedContact(ctx, principalType, principal)
	if err != nil {
		logger.Error("userRepo.GetUserAndContactByConfirmedContact call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.User{}, err
	}
	span.AddEvent("user and contact retreived from repo")
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && now.Before(user.LockedOutUntil.Value) {
		err := coreerrors.NewUserLockedOutError(user.ID, true)
		logger.Error(err.GetErrorMessage(), zap.Any("error", err))
		evtString := fmt.Sprintf("user is locked out until %s", user.LockedOutUntil.Value.UTC().String())
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, err
	}
	if !contact.IsPrimary {
		err := coreerrors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
		logger.Error(err.GetErrorMessage(), zap.Any("error", err))
		evtString := fmt.Sprintf("contact user is not primary: %s of type %s", contact.Principal, contact.Type)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, err
	}
	// TODO: remove this I guess because if the contact is not confirmed then GetUserAndContactByConfirmedContact will not return anything
	if !contact.IsConfirmed() { // || contact.ConfirmedDate.Value.After(now)
		err := coreerrors.NewLoginPrimaryContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
		logger.Error(err.GetErrorMessage(), zap.Any("error", err))
		evtString := fmt.Sprintf("contact is not confirmed: %s of type %s", contact.Principal, contact.Type)
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, err
	}
	span.AddEvent("contact validated")
	// check password
	// hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	passwordMatch, err := utilities.BcryptCompareStringAndHash(user.PasswordHash, password, user.ID)
	if err != nil {
		evtString := "failed to check users password hash"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
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
			logger.Error("update user after consecutive failed login increment failed", zap.Any("error", err))
		}
		err = coreerrors.NewLoginFailedWrongPasswordError(user.ID, true)
		evtString := err.GetErrorMessage()
		logger.Warn(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, err
	}
	user.LastLoginDate.Set(now)
	user.ConsecutiveFailedLoginAttempts = 0
	user.LockedOutUntil.Unset()
	err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
	if err != nil {
		evtString := "update user after successful login"
		apptelemetry.SetSpanError(&span, err, evtString)
		logger.Error(evtString, zap.Any("error", err))
		return models.User{}, err
	}
	span.AddEvent("login completed")
	return user, nil
}

// TODO: remove string from return and make work like rgistration call. test with stackemailservice
func (ls loginService) StartPasswordResetByPrimaryContact(ctx context.Context, logger *zap.Logger, principal, principalType string, initiator string) (string, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, ls.GetName(), "StartPasswordResetByPrimaryContact")
	defer span.End()
	user, contact, err := ls.userRepo.GetUserAndContactByConfirmedContact(ctx, principalType, principal)
	if err != nil {
		logger.Error("userRepo.GetUserAndContactByContact call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return "", err
	}
	span.AddEvent("user and contact retreived from repo")
	if !contact.IsPrimary {
		evtString := fmt.Sprintf("contact user is not primary: %s of type %s", contact.Principal, contact.Type)
		err := coreerrors.NewPasswordResetContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return "", err
	}
	// TODO: make password reset token expiration configurable.
	token, err := models.NewToken(user.ID, models.TokenTypePasswordReset, time.Minute*15)
	if err != nil {
		evtString := "failed to create new password reset token"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return "", err
	}
	span.AddEvent("new password reset token created")
	err = ls.tokenService.PutToken(ctx, logger, token)
	if err != nil {
		evtString := "failed to store new password reset token"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return "", err
	}
	span.AddEvent("new password reset token stored in repo")
	switch contact.Type {
	case core.CONTACT_TYPE_EMAIL:
		// TODO: create template for this...
		body := fmt.Sprintf("A Password reset has been initiated. Your password reset token is: %s", token.Value)
		err = ls.emailService.SendPlainTextEmail(ctx, logger, []string{contact.Principal}, "Password reset", body)
		if err != nil {
			evtString := "failed to send password reset notification error occurred"
			logger.Error(evtString, zap.Any("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return token.Value, err // TODO: what should we do here???
		}
	default:
		err := coreerrors.NewComponentNotImplementedError("notification system", fmt.Sprintf("%s notification service", contact.Type), true)
		evtString := fmt.Sprintf("failed to send notification contact type not supported: %s", contact.Type)
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return "", err
	}
	span.AddEvent("password reset initiated")
	return token.Value, nil
}

func (ls loginService) ResetPassword(ctx context.Context, logger *zap.Logger, passwordResetToken string, newPassword string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, ls.GetName(), "ResetPassword")
	defer span.End()
	if newPassword == "" {
		err := coreerrors.NewNoNewPasswordHashProvidedError(true)
		evtString := "new password is empty string"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// TODO: add password validation logic
	token, err := ls.tokenService.GetToken(ctx, logger, passwordResetToken, models.TokenTypePasswordReset)
	if err != nil {
		logger.Error("tokenService.GetToken call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("password reset token retreived from repo")
	user, err := ls.userRepo.GetUserByID(ctx, token.TargetID)
	if err != nil {
		logger.Error("userRepo.GetUserByID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("user retreived from repo")
	newPasswordHash, err := utilities.BcryptHashString(newPassword, bcrypt.DefaultCost)
	if err != nil {
		evtString := "failed to hash users new password"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("new password hash generated")
	user.PasswordHash = newPasswordHash
	err = ls.userRepo.UpdateUser(ctx, &user, initiator)
	if err != nil {
		logger.Error("userRepo.UpdateUser call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("password reset completed")
	return nil
}
