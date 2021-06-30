package core

import "context"

// TODO: implemet account lockout after X number of consecutive failed login attempts.
type LoginService interface {
	// LoginWithContact attempts to confirm a users credentials and if they match it returns true and resets the users ConsecutiveFailedLoginAttempts, otherwise it returns false and increments the users ConsecutiveFailedLoginAttempts
	LoginWithContact(ctx context.Context, principal, principalType, password string) (bool, error)
	// StartPasswordResetByContact sets a password reset token for the user with the corresponding principal and type that are confirmed.
	StartPasswordResetByContact(ctx context.Context, principal, principalType string) (string, error)
}

type UserDataService interface {
}
