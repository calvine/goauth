package service

import (
	"context"
	"fmt"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

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

func (us userService) GetUserAndContactByConfirmedContact(ctx context.Context, logger *zap.Logger, contactType string, contactPrincipal string, initiator string) (models.User, models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUserAndContactByConfirmedContact")
	defer span.End()
	user, contact, err := us.userRepo.GetUserAndContactByConfirmedContact(ctx, contactType, contactPrincipal)
	if err != nil {
		logger.Error("userRepo.GetUserAndContactByContact call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.User{}, models.Contact{}, err
	}
	if !contact.IsConfirmed() {
		evtString := fmt.Sprintf("contact found is not confirmed: ID = %s", contact.ID)
		err := coreerrors.NewRegisteredContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, models.Contact{}, err
	}
	return user, contact, nil
}

func (us userService) RegisterUserAndPrimaryContact(ctx context.Context, logger *zap.Logger, contactType, contactPrincipal string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "RegisterUserAndPrimaryContact")
	defer span.End()
	// check that email address does not already exist as a confirmed contact.
	numExistingConfirmedContacts, err := us.contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType(ctx, contactType, contactPrincipal)
	if err != nil {
		logger.Error("contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	if numExistingConfirmedContacts != 0 {
		var err errors.RichError
		if numExistingConfirmedContacts > 1 {
			// This is really bad and we need to know about it asap!
			errMsg := "critical issue here more than one contact is confirmed with this info"
			err = coreerrors.NewMultipleConfirmedInstancesOfContactError(contactPrincipal, contactType, numExistingConfirmedContacts, true)
			logger.Error(errMsg, zap.Any("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, errMsg)
		} else {
			errMsg := "a contact already exists and is confirmed with the data provided"
			errorFields := make(map[string]interface{})
			errorFields["numExistingConfirmedContacts"] = numExistingConfirmedContacts
			err = coreerrors.NewRegistrationContactAlreadyConfirmedError(contactPrincipal, contactType, errorFields, true)
			logger.Error(errMsg,
				zap.Any("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, errMsg)
		}
		return err
	}
	// create new user and contact in datastore

	// send confirmation email

	// NOTE: allow user to set password on confirmation link click.

	return coreerrors.NewNotImplementedError(true)
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
