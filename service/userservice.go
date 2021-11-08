package service

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
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
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUserByConfirmedContact")
	defer span.End()
	return models.User{}, coreerrors.NewNotImplementedError(true)
}

func (us userService) RegisterUserAndPrimaryContact(ctx context.Context, logger *zap.Logger, contactPrincipal, contactType string) (models.User, models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "RegisterUserAndPrimaryContact")
	defer span.End()
	// check to see if a user is registered with the confirmed contact info provided
	// existingUser, existingContact, err := us.userRepo.GetUserAndContactByContact(ctx, contactType, contactPrincipal)
	// if err != nil {
	// 	logger.Error("", zap.Error(err))
	// }
	return models.User{}, models.Contact{}, coreerrors.NewNotImplementedError(true)
}

// func (us userService) AddUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError {
// 	return coreerrors.NewNotImplementedError(true)
// }

// func (us userService) UpdateUser(ctx context.Context, logger *zap.Logger, user *models.User, initiator string) errors.RichError {
// 	return coreerrors.NewNotImplementedError(true)
// }

func (us userService) GetUserPrimaryContact(ctx context.Context, logger *zap.Logger, userID string, initiator string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUserPrimaryContact")
	defer span.End()
	return models.Contact{}, coreerrors.NewNotImplementedError(true)
}

func (us userService) GetUsersContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersContacts")
	defer span.End()
	return nil, coreerrors.NewNotImplementedError(true)
}

func (us userService) GetUsersConfirmedContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersConfirmedContacts")
	defer span.End()
	return nil, coreerrors.NewNotImplementedError(true)
}

func (us userService) AddContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "AddContact")
	defer span.End()
	return coreerrors.NewNotImplementedError(true)
}

func (us userService) UpdateContact(ctx context.Context, logger *zap.Logger, contact *models.Contact, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "UpdateContact")
	defer span.End()
	return coreerrors.NewNotImplementedError(true)
}

func (us userService) ConfirmContact(ctx context.Context, logger *zap.Logger, confirmationCode string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "ConfirmContact")
	defer span.End()
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
