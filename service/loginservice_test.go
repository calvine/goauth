package service

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core/constants/contact"
	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/goauth/dataaccess/memory"
	"go.uber.org/zap/zaptest"
	"golang.org/x/crypto/bcrypt"
)

var (
	loginServiceTest_ConfirmedUser      models.User
	loginServiceTest_UnconfirmedUser    models.User
	loginServiceTest_OtherConfirmedUser models.User

	loginServiceTest__ConfirmedPrimaryContact     models.Contact
	loginServiceTest_ConfirmedSecondaryContact    models.Contact
	loginServiceTest_UnconfirmedPrimaryContact    models.Contact
	loginServiceTest_OtherConfirmedPrimaryContact models.Contact

	loginServiceTest_TestPasswordResetToken string

	loginServiceTest_NonPasswordResetToken models.Token
)

const (
	loginServiceTest_CreatedBy = "login service tests"

	loginServiceTest_ConfirmedPrimaryEmail      = "confirmed@email.com"
	loginServiceTest_ConfirmedSecondaryEmail    = "secondary@email.com"
	loginServiceTest_UnconfirmedPrimaryEmail    = "unconfirmed@email.com"
	loginServiceTest_OtherConfirmedPrimaryEmail = "otherconfirmed@email.com"

	loginServiceTest_ConfirmedUserPassword      = "testpass"
	loginServiceTest_UnconfirmedUserPassword    = "tp2"
	loginServiceTest_OtherConfirmedUserPassword = "testpass3"

	loginServiceTest_NewPasswordPostReset = "anewpasswordhash123"

	loginServiceTest_LockoutAfterFailedLoginAttempts = 10

	loginServiceTest_LockoutDuration time.Duration = time.Millisecond * 500

	loginServiceTest_LockoutReleaseWaitDuration time.Duration = time.Millisecond * 700
)

func TestLoginService(t *testing.T) {
	loginService := buildLoginService(t)

	t.Run("GetName", func(t *testing.T) {
		_testLoginServiceGetName(t, loginService)
	})

	t.Run("StartPasswordResetByPrimaryContact", func(t *testing.T) {
		_testStartPasswordResetByPrimaryContact(t, loginService)
	})

	t.Run("ResetPassword", func(t *testing.T) {
		_testResetPassword(t, loginService)
	})

	t.Run("LoginWithPrimaryContact", func(t *testing.T) {
		_testLoginWithPrimaryContact(t, loginService)
	})
}

func setupLoginServiceTestData(t *testing.T, userRepo repo.UserRepo, contactRepo repo.ContactRepo, tokenService services.TokenService) {
	passHash, err := utilities.BcryptHashString(loginServiceTest_ConfirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create test password hash: %s", err.GetErrorCode())
		t.FailNow()
	}
	passHash2, err := utilities.BcryptHashString(loginServiceTest_UnconfirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create test password hash: %s", err.GetErrorCode())
		t.FailNow()
	}
	passHash3, err := utilities.BcryptHashString(loginServiceTest_OtherConfirmedUserPassword, bcrypt.DefaultCost)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to create test password hash: %s", err.GetErrorCode())
		t.FailNow()
	}
	loginServiceTest_ConfirmedUser = models.User{
		ID:           "123",
		PasswordHash: passHash,
	}
	loginServiceTest__ConfirmedPrimaryContact = models.Contact{
		ID:            "456",
		UserID:        loginServiceTest_ConfirmedUser.ID,
		Type:          contact.Email,
		Principal:     loginServiceTest_ConfirmedPrimaryEmail,
		IsPrimary:     true,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}
	loginServiceTest_ConfirmedSecondaryContact = models.Contact{
		ID:            "789",
		UserID:        loginServiceTest_ConfirmedUser.ID,
		Type:          contact.Email,
		Principal:     loginServiceTest_ConfirmedSecondaryEmail,
		IsPrimary:     false,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}

	loginServiceTest_UnconfirmedUser = models.User{
		ID:           "012",
		PasswordHash: passHash2,
	}
	loginServiceTest_UnconfirmedPrimaryContact = models.Contact{
		ID:        "345",
		UserID:    loginServiceTest_UnconfirmedUser.ID,
		Principal: loginServiceTest_UnconfirmedPrimaryEmail,
		IsPrimary: true,
		Type:      contact.Email,
	}

	loginServiceTest_OtherConfirmedUser = models.User{
		ID:           "678",
		PasswordHash: passHash3,
	}
	loginServiceTest_OtherConfirmedPrimaryContact = models.Contact{
		ID:            "901",
		UserID:        loginServiceTest_OtherConfirmedUser.ID,
		Type:          contact.Email,
		Principal:     loginServiceTest_OtherConfirmedPrimaryEmail,
		IsPrimary:     true,
		ConfirmedDate: nullable.NullableTime{HasValue: true, Value: time.Now().Add(time.Minute * -1)},
	}

	err = userRepo.AddUser(context.TODO(), &loginServiceTest_ConfirmedUser, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add user for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = userRepo.AddUser(context.TODO(), &loginServiceTest_UnconfirmedUser, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add user for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = userRepo.AddUser(context.TODO(), &loginServiceTest_OtherConfirmedUser, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add user for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = contactRepo.AddContact(context.TODO(), &loginServiceTest__ConfirmedPrimaryContact, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add contact for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = contactRepo.AddContact(context.TODO(), &loginServiceTest_ConfirmedSecondaryContact, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add contact for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = contactRepo.AddContact(context.TODO(), &loginServiceTest_UnconfirmedPrimaryContact, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add contact for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	err = contactRepo.AddContact(context.TODO(), &loginServiceTest_OtherConfirmedPrimaryContact, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add contact for login service tests: %s", err.GetErrorCode())
		t.FailNow()
	}

	loginServiceTest_NonPasswordResetToken, err = models.NewToken("", models.TokenTypeSession, time.Minute*10)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add non password reset token with error: %s", err.GetErrorCode())
		t.FailNow()
	}
	logger := zaptest.NewLogger(t)
	tokenService.PutToken(context.TODO(), logger, loginServiceTest_NonPasswordResetToken)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("failed to add non password reset token with error: %s", err.GetErrorCode())
		t.FailNow()
	}
}

func buildLoginService(t *testing.T) services.LoginService {
	auditLogRepo := memory.NewMemoryAuditLogRepo(false)
	users := make(map[string]models.User)
	contacts := make(map[string]models.Contact)
	userRepo, err := memory.NewMemoryUserRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	contactRepo, err := memory.NewMemoryContactRepo(&users, &contacts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	tokenRepo := memory.NewMemoryTokenRepo()
	emailService, _ := NewEmailService(NoOpEmailService, nil)
	tokenService := NewTokenService(tokenRepo)

	setupLoginServiceTestData(t, userRepo, contactRepo, tokenService)

	options := LoginServiceOptions{
		AuditLogRepo:           auditLogRepo,
		ContactRepo:            contactRepo,
		UserRepo:               userRepo,
		EmailService:           emailService,
		TokenService:           tokenService,
		MaxFailedLoginAttempts: loginServiceTest_LockoutAfterFailedLoginAttempts,
		AccountLockoutDuration: loginServiceTest_LockoutDuration,
	}

	return NewLoginService(options)
}

func _testLoginServiceGetName(t *testing.T, loginService services.LoginService) {
	serviceName := loginService.GetName()
	expectedServiceName := "loginService"
	if serviceName != expectedServiceName {
		t.Errorf("service name is not what was expected: got %s - expected %s", serviceName, expectedServiceName)
	}
}

func _testStartPasswordResetByPrimaryContact(t *testing.T, loginService services.LoginService) {
	// successful start password reset request
	t.Run("Success", func(t *testing.T) {
		__testStartPasswordResetSuccess(t, loginService)
	})

	// failed start password reset with non primary contact
	t.Run("Failed not primary contact", func(t *testing.T) {
		__testStartPasswordResetFailedNotPrimaryContact(t, loginService)
	})
}

func __testStartPasswordResetSuccess(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	tokenValue, err := loginService.StartPasswordResetByPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_ConfirmedPrimaryEmail, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("received error when attempting to start valid password reset: %s", err.GetErrorCode())
	}
	if tokenValue == "" {
		t.Error("token value should not be an empty string")
	}
	loginServiceTest_TestPasswordResetToken = tokenValue
}

func __testStartPasswordResetFailedNotPrimaryContact(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	tokenValue, err := loginService.StartPasswordResetByPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_ConfirmedSecondaryEmail, loginServiceTest_CreatedBy)
	if err == nil {
		t.Error("expected error due to non primary contact being used for password reset")
	}
	if err.GetErrorCode() != errors.ErrCodePasswordResetContactNotPrimary {
		t.Log(err.Error())
		t.Errorf("expected password reset contact not primary error but got: %s", err.GetErrorCode())
	}
	if tokenValue != "" {
		t.Errorf("token value should be empty because the password reset initiation should have failed: %s", tokenValue)
	}
}

func _testResetPassword(t *testing.T, loginService services.LoginService) {
	// password reset successful
	t.Run("Success", func(t *testing.T) {
		__testPasswordResetSuccess(t, loginService)
	})

	// password reset failure invalid token
	t.Run("Failure invalid token", func(t *testing.T) {
		__testPasswordResetFailureInvalidToken(t, loginService)
	})

	// password reset failure empty password hash
	t.Run("Failure empty password hash", func(t *testing.T) {
		__testPasswordResetFailureEmptyPasswordHash(t, loginService)
	})

	// password reset failure non password reset token presented
	t.Run("Failure wrong token type", func(t *testing.T) {
		__testPasswordResetFailureWrongTokenType(t, loginService)
	})
}

func __testPasswordResetSuccess(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	err := loginService.ResetPassword(context.TODO(), logger, loginServiceTest_TestPasswordResetToken, loginServiceTest_NewPasswordPostReset, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("expected password reset to succeed bug got an an error: %s", err.GetErrorCode())
	}
}

func __testPasswordResetFailureInvalidToken(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	err := loginService.ResetPassword(context.TODO(), logger, "made up token that is not real", "new password hash 2", loginServiceTest_CreatedBy)
	if err == nil {
		t.Errorf("expected password reset to fail because the token provided is not a real token")
	}
}

func __testPasswordResetFailureEmptyPasswordHash(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	err := loginService.ResetPassword(context.TODO(), logger, "made up token that is not real", "", loginServiceTest_CreatedBy)
	if err == nil {
		t.Errorf("expected password reset to fail because the token provided is not a real token")
	}
	if err.GetErrorCode() != errors.ErrCodeNoNewPasswordHashProvided {
		t.Log(err.Error())
		t.Errorf("expected no new password hash error but got: %s", err.GetErrorCode())
	}
}

func __testPasswordResetFailureWrongTokenType(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	err := loginService.ResetPassword(context.TODO(), logger, loginServiceTest_NonPasswordResetToken.Value, "new password hash 3", loginServiceTest_CreatedBy)
	if err == nil {
		t.Errorf("expected password reset to fail because the token provided is not a password reset token")
	}
	if err.GetErrorCode() != errors.ErrCodeWrongTokenType {
		t.Log(err.Error())
		t.Errorf("expected no new password hash error but got: %s", err.GetErrorCode())
	}
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
	logger := zaptest.NewLogger(t)
	user, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_ConfirmedPrimaryEmail, loginServiceTest_NewPasswordPostReset, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("confirmed user login with primary contact email should be successfull: %s", err.GetErrorCode())
	}
	if loginServiceTest_ConfirmedUser.ID != user.ID {
		t.Errorf("expected user id does not match returned user id: got %s - expected %s", loginServiceTest_ConfirmedUser.ID, user.ID)
	}
	if !user.LastLoginDate.HasValue {
		t.Error("user.LastLoginDate should be set.")
	}
}

func __testFailedLogin(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_ConfirmedPrimaryEmail, "not the right password 12345678904321234567", loginServiceTest_CreatedBy)
	if err == nil {
		t.Error("expected failed login wrong password error bug got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeLoginFailedWrongPassword {
		t.Log(err.Error())
		t.Errorf("expected failed login wrong password error bug got another error: %s", err.GetErrorCode())
	}
}

func __testFailedLoginPrimaryContactNotConfirmed(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_UnconfirmedPrimaryEmail, "not the right password 12345678904321234567", loginServiceTest_CreatedBy)
	if err == nil {
		t.Error("expected failed login no user found got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeNoUserFound {
		t.Log(err.Error())
		t.Errorf("expected failed login no user found another error: %s", err.GetErrorCode())
	}
}

func __testFailedLoginSecondaryContactUsed(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_ConfirmedSecondaryEmail, loginServiceTest_ConfirmedUserPassword, loginServiceTest_CreatedBy)
	if err == nil {
		t.Error("expected failed login wrong password error bug got no error")
	}
	if err.GetErrorCode() != errors.ErrCodeLoginContactNotPrimary {
		t.Log(err.Error())
		t.Errorf("expected failed login wrong password error bug got another error: %s", err.GetErrorCode())
	}
}

func __testAccountLockout(t *testing.T, loginService services.LoginService) {
	logger := zaptest.NewLogger(t)
	for i := 0; i < loginServiceTest_LockoutAfterFailedLoginAttempts; i++ {
		_, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_OtherConfirmedPrimaryEmail, "Not the right password34567898765trew2123456&*!&^", loginServiceTest_CreatedBy)
		if err == nil {
			t.Error("expected error because login details are invalid...")
		}
		if err.GetErrorCode() != errors.ErrCodeLoginFailedWrongPassword {
			t.Log(err.Error())
			t.Errorf("expected failed login wrong password error but got error of type: %s", err.GetErrorCode())
		}
	}
	_, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_OtherConfirmedPrimaryEmail, "Not the right password34567898765trew2123456&*!&^", loginServiceTest_CreatedBy)
	if err == nil {
		t.Errorf("should have locked out user account after %d failed login attempts", loginServiceTest_LockoutAfterFailedLoginAttempts)
	}
	if err.GetErrorCode() != errors.ErrCodeUserLockedOut {
		t.Log(err.Error())
		t.Errorf("expected user locked out error but received: %s", err.GetErrorCode())
	}
}

func __testAccountLockoutRelease(t *testing.T, loginService services.LoginService) {
	time.Sleep(loginServiceTest_LockoutReleaseWaitDuration)
	logger := zaptest.NewLogger(t)
	user, err := loginService.LoginWithPrimaryContact(context.TODO(), logger, contact.Email, loginServiceTest_OtherConfirmedPrimaryEmail, loginServiceTest_OtherConfirmedUserPassword, loginServiceTest_CreatedBy)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("after sleep the users account lockout should have expired but got this error: %s", err.GetErrorCode())
	}
	if user.ConsecutiveFailedLoginAttempts != 0 {
		t.Errorf("expected user consecutive failed login attempts to be 0 but got %d", user.ConsecutiveFailedLoginAttempts)
	}
	if user.LockedOutUntil.HasValue {
		t.Errorf("expected user locked out until to be unset but got : %s", user.LockedOutUntil.Value.String())
	}
}
