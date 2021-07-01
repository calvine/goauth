package service

import (
	repo "github.com/calvine/goauth/core/repositories"
)

type loginService struct {
	userRepo     repo.UserRepo
	contactRepo  repo.ContactRepo
	auditLogRepo repo.AuditLogRepo
}

func NewLoginService(userRepo repo.UserRepo, contactRepo repo.ContactRepo, auditLogRepo repo.AuditLogRepo) *loginService {
	return &loginService{
		userRepo:     userRepo,
		contactRepo:  contactRepo,
		auditLogRepo: auditLogRepo,
	}
}

// func (ls *loginService) LoginWithContact(ctx context.Context, principal, principalType, password string, initiator string) (bool, error) {

// }

// func (ls *loginService) StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, error) {

// }

// func (ls *loginService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, error) {
// 	contact, err := ls.contactRepo.GetContactByConfirmationCode(ctx, confirmationCode)
// 	return false, nil
// }
