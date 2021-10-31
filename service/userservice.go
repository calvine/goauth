package service

import (
	"context"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

// TODO: need to refine service methods... like the add user method should also take a proposed primary contact...

type userService struct {
	userRepo     repo.UserRepo
	contactRepo  repo.ContactRepo
	tokenService services.TokenService
	emailService services.EmailService
}

func NewUserService(userRepo repo.UserRepo, contactRepo repo.ContactRepo, tokenService services.TokenService, emailService services.EmailService) services.UserService {
	return userService{
		userRepo:     userRepo,
		contactRepo:  contactRepo,
		tokenService: tokenService,
		emailService: emailService,
	}
}

func (us userService) GetName() string {
	return "userService"
}

func (us userService) GetUserByConfirmedContact(ctx context.Context, logger *zap.Logger, contactPrincipal string, initiator string) (models.User, errors.RichError) {
	return models.User{}, coreerrors.NewNotImplementedError(true)
}

// func (us userService) AddUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError {
// 	return coreerrors.NewNotImplementedError(true)
// }

// func (us userService) UpdateUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError {
// 	return coreerrors.NewNotImplementedError(true)
// }

func (us userService) GetUserPrimaryContact(ctx context.Context, logger *zap.Logger, userID string, initiator string) (models.Contact, errors.RichError) {
	return models.Contact{}, coreerrors.NewNotImplementedError(true)
}

func (us userService) GetUsersContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	return nil, coreerrors.NewNotImplementedError(true)
}

func (us userService) GetUsersConfirmedContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	return nil, coreerrors.NewNotImplementedError(true)
}

func (us userService) AddContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (us userService) UpdateContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (us userService) ConfirmContact(ctx context.Context, logger *zap.Logger, confirmationCode string, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

// func (us userService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, errors.RichError) {
// 	token, err := ls.tokenService.GetToken(confirmationCode, models.TokenTypeConfirmContact)
// 	if err != nil {
// 		return false, err
// 	}
// 	contact, err := ls.contactRepo.GetContactByContactID(ctx, token.TargetID)
// 	if err != nil {
// 		return false, err
// 	}
// 	now := time.Now().UTC()
// 	contact.ConfirmedDate.Set(now)
// 	err = ls.contactRepo.UpdateContact(ctx, &contact, initiator)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
