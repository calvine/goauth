package service

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/goauth/dataaccess/memory"
	"golang.org/x/crypto/bcrypt"
)

var (
	confirmedUser      models.User
	unconfirmedUser    models.User
	otherConfirmedUser models.User

	confirmedPrimaryContact      models.Contact
	confirmedSecondaryContact    models.Contact
	unconfirmedPrimaryContact    models.Contact
	otherConfirmedPrimaryContact models.Contact
)

const (
	loginServiceTestCreatedBy = "login service tests"

	confirmedPrimaryEmail      = "confirmed@email.com"
	confirmedSecondaryEmail    = "secondary@email.com"
	unconfirmedPrimaryEmail    = "unconfirmed@email.com"
	otherConfirmedPrimaryEmail = "otherconfirmed@email.com"

	confirmedUserPassword      = "testpass"
	unconfirmedUserPassword    = "tp2"
	otherConfirmedUserPassword = "testpass3"

	lockoutAfterFailedLoginAttempts = 10

	lockoutDuration time.Duration = time.Millisecond * 500

	lockoutReleaseWaitDuration time.Duration = time.Millisecond * 700
)

func TestLoginService(t *testing.T) {
	loginService := buildLoginService(t)

	t.Run("LoginWithPrimaryContact", func(t *testing.T) {
		_testLoginWithPrimaryContact(t, loginService)
	})

	t.Run("StartPasswordResetByContact", func(t *testing.T) {
		_testStartPasswordResetByContact(t, loginService)
	})

	t.Run("ResetPassword", func(t *testing.T) {
		_testResetPassword(t, loginService)
	})
}

func setupTestData(t *testing.T, userRepo repo.UserRepo, contactRepo repo.ContactRepo) {
	passHash, err := utilities.BcryptHashString(confirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("failed to create test password hash: %s", err.Error())
	}
	passHash2, err := utilities.BcryptHashString(unconfirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("failed to create test password hash: %s", err.Error())
	}
	passHash3, err := utilities.BcryptHashString(otherConfirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("failed to create test password hash: %s", err.Error())
	}
	confirmedUser = models.User{
		ID:           "123",
		PasswordHash: passHash,
	}
	confirmedPrimaryContact = models.Contact{
		ID:            "456",
		UserID:        confirmedUser.ID,
		Type:          core.CONTACT_TYPE_EMAIL,
		Principal:     confirmedPrimaryEmail,
		IsPrimary:     true,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}
	confirmedSecondaryContact = models.Contact{
		ID:            "789",
		UserID:        confirmedUser.ID,
		Type:          core.CONTACT_TYPE_EMAIL,
		Principal:     confirmedSecondaryEmail,
		IsPrimary:     false,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}

	unconfirmedUser = models.User{
		ID:           "012",
		PasswordHash: passHash2,
	}
	unconfirmedPrimaryContact = models.Contact{
		ID:        "345",
		UserID:    unconfirmedUser.ID,
		Principal: unconfirmedPrimaryEmail,
		IsPrimary: true,
		Type:      core.CONTACT_TYPE_EMAIL,
	}

	otherConfirmedUser = models.User{
		ID:           "678",
		PasswordHash: passHash3,
	}
	otherConfirmedPrimaryContact = models.Contact{
		ID:            "901",
		UserID:        otherConfirmedUser.ID,
		Type:          core.CONTACT_TYPE_EMAIL,
		Principal:     otherConfirmedPrimaryEmail,
		IsPrimary:     true,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}

	err = userRepo.AddUser(context.TODO(), &confirmedUser, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add user for login service tests: %s", err.Error())
	}

	err = userRepo.AddUser(context.TODO(), &unconfirmedUser, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add user for login service tests: %s", err.Error())
	}

	err = userRepo.AddUser(context.TODO(), &otherConfirmedUser, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add user for login service tests: %s", err.Error())
	}

	err = contactRepo.AddContact(context.TODO(), &confirmedPrimaryContact, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add contact for login service tests: %s", err.Error())
	}

	err = contactRepo.AddContact(context.TODO(), &confirmedSecondaryContact, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add contact for login service tests: %s", err.Error())
	}

	err = contactRepo.AddContact(context.TODO(), &unconfirmedPrimaryContact, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add contact for login service tests: %s", err.Error())
	}

	err = contactRepo.AddContact(context.TODO(), &otherConfirmedPrimaryContact, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("failed to add contact for login service tests: %s", err.Error())
	}
}

func buildLoginService(t *testing.T) services.LoginService {
	auditLogRepo := memory.NewMemoryAuditLogRepo(false)
	userRepo := memory.NewMemoryUserRepo()
	contactRepo := memory.NewMemoryContactRepo()
	tokenRepo := memory.NewMemoryTokenRepo()
	emailService, _ := NewEmailService(MockEmailService, nil)
	tokenService := NewTokenService(tokenRepo)

	setupTestData(t, userRepo, contactRepo)

	options := LoginServiceOptions{
		AuditLogRepo:           auditLogRepo,
		ContactRepo:            contactRepo,
		UserRepo:               userRepo,
		EmailService:           emailService,
		TokenService:           tokenService,
		MaxFailedLoginAttempts: lockoutAfterFailedLoginAttempts,
		AccountLockoutDuration: lockoutDuration,
	}

	return NewLoginService(options)
}

func _testLoginWithPrimaryContact(t *testing.T, loginService services.LoginService) {
	// test successfull login
	t.Run("Successfull email login", func(t *testing.T) {
		__testSuccessfullEmailLogin(t, loginService)
	})

	// test falied login
	t.Run("Failed email login password", func(t *testing.T) {
		__testFailedLogin(t, loginService)
	})

	// test login failed primary contact not confirmed
	t.Run("Failed email login primary contact not confirmed", func(t *testing.T) {
		__testFailedLoginPrimaryContactNotConfirmed(t, loginService)
	})

	// test login failed secondary contact used
	t.Run("Failed email login secondary contact used", func(t *testing.T) {
		__testFailedLoginSecondaryContactUsed(t, loginService)
	})

	// test account lockout
	t.Run("Account lockout", func(t *testing.T) {
		__testAccountLockout(t, loginService)
	})

	// test account lockout release?
	t.Run("Account lockout release", func(t *testing.T) {
		__testAccountLockoutRelease(t, loginService)
	})
}

func __testSuccessfullEmailLogin(t *testing.T, loginService services.LoginService) {
	user, err := loginService.LoginWithPrimaryContact(context.TODO(), confirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, confirmedUserPassword, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("confirmed user login with primary contact email should be successfull: %s", err.Error())
	}
	if confirmedUser.ID != user.ID {
		t.Errorf("expected user id does not match returned user id: got %s - expected %s", confirmedUser.ID, user.ID)
	}
	if !user.LastLoginDate.HasValue {
		t.Error("user.LastLoginDate should be set.")
	}
}

func __testFailedLogin(t *testing.T, loginService services.LoginService) {
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), confirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, "not the right password 12345678904321234567", loginServiceTestCreatedBy)
	if err == nil {
		t.Error("expected failed login wrong password error bug got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeLoginFailedWrongPassword {
		t.Errorf("expected failed login wrong password error bug got another error: %s - %s", err.GetErrorCode(), err.Error())
	}
}

func __testFailedLoginPrimaryContactNotConfirmed(t *testing.T, loginService services.LoginService) {
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), unconfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, "not the right password 12345678904321234567", loginServiceTestCreatedBy)
	if err == nil {
		t.Error("expected failed login wrong password error bug got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeContactNotConfirmed {
		t.Errorf("expected failed login wrong password error bug got another error: %s - %s", err.GetErrorCode(), err.Error())
	}
}

func __testFailedLoginSecondaryContactUsed(t *testing.T, loginService services.LoginService) {
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), confirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, confirmedUserPassword, loginServiceTestCreatedBy)
	if err == nil {
		t.Error("expected failed login wrong password error bug got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeLoginContactNotPrimary {
		t.Errorf("expected failed login wrong password error bug got another error: %s - %s", err.GetErrorCode(), err.Error())
	}
}

func __testAccountLockout(t *testing.T, loginService services.LoginService) {
	for i := 0; i < lockoutAfterFailedLoginAttempts; i++ {
		_, err := loginService.LoginWithPrimaryContact(context.TODO(), otherConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, "Not the right password34567898765trew2123456&*!&^", loginServiceTestCreatedBy)
		if err == nil {
			t.Error("expected error because login details are invalid...")
		}
		if err.GetErrorCode() != errors.ErrCodeLoginFailedWrongPassword {
			t.Errorf("expected failed login wrong password error but got error of type %s: %s", err.GetErrorCode(), err.Error())
		}
	}
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), otherConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, "Not the right password34567898765trew2123456&*!&^", loginServiceTestCreatedBy)
	if err == nil {
		t.Errorf("should have locked out user account after %d failed login attempts", lockoutAfterFailedLoginAttempts)
	}
	if err.GetErrorCode() != errors.ErrCodeUserLockedOut {
		t.Errorf("expected user locked out error but received: %s", err.Error())
	}
}

func __testAccountLockoutRelease(t *testing.T, loginService services.LoginService) {
	time.Sleep(lockoutReleaseWaitDuration)
	user, err := loginService.LoginWithPrimaryContact(context.TODO(), otherConfirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, otherConfirmedUserPassword, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("after sleep the users account lockout should have expired but got this error: %s", err.Error())
	}
	if user.ConsecutiveFailedLoginAttempts != 0 {
		t.Errorf("expected user consecutive failed login attempts to be 0 but got %d", user.ConsecutiveFailedLoginAttempts)
	}
	if user.LockedOutUntil.HasValue {
		t.Errorf("expected user locked out until to be unset but got : %s", user.LockedOutUntil.Value.String())
	}
}

func _testStartPasswordResetByContact(t *testing.T, loginService services.LoginService) {
	t.Error("test not implemented")
	// successful password reset request
	t.Run("StartPasswordResetSuccess", func(t *testing.T) {
		__testStartPasswordResetSuccess(t, loginService)
	})

	// failed password reset with non primary contact
	t.Run("StartPasswordResetFailedNotPrimaryContact", func(t *testing.T) {
		__testStartPasswordResetFailedNotPrimaryContact(t, loginService)
	})
}

func __testStartPasswordResetSuccess(t *testing.T, loginService services.LoginService) {
	tokenValue, err := loginService.StartPasswordResetByContact(context.TODO(), confirmedPrimaryEmail, core.CONTACT_TYPE_EMAIL, loginServiceTestCreatedBy)
	if err != nil {
		t.Errorf("received error when attempting to start valid password reset: %s", err.Error())
	}
	if tokenValue == "" {
		t.Error("token value should not be an empty string")
	}
}

func __testStartPasswordResetFailedNotPrimaryContact(t *testing.T, loginService services.LoginService) {
	tokenValue, err := loginService.StartPasswordResetByContact(context.TODO(), confirmedSecondaryEmail, core.CONTACT_TYPE_EMAIL, loginServiceTestCreatedBy)
	if err == nil {
		t.Error("expected error due to non primary contact being used for password reset")
	}
	if err.GetErrorCode() != errors.ErrCodePasswordResetContactNotPrimary {
		t.Errorf("expected password reset contact not primary error but got %s: %s", err.GetErrorCode(), err.Error())
	}
	if tokenValue != "" {
		t.Errorf("token value should be empty because the password reset initiation should have failed: %s", tokenValue)
	}
}

func _testResetPassword(t *testing.T, loginService services.LoginService) {
	t.Error("test not implemented")
}
