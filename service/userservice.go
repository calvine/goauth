package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
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
	span.AddEvent("user and contact retreived")
	return user, contact, nil
}

func (us userService) RegisterUserAndPrimaryContact(ctx context.Context, logger *zap.Logger, contactType, contactPrincipal string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "RegisterUserAndPrimaryContact")
	defer span.End()
	// TODO: normalize contact principal
	// check that email address does not already exist as a confirmed contact.
	err := us.checkForExistingConfirmedContacts(ctx, logger, &span, contactType, contactPrincipal, "")
	if err != nil {
		return err
	}
	// create new user and contact in datastore
	newUser := models.NewUser()
	err = us.userRepo.AddUser(ctx, &newUser, initiator)
	if err != nil {
		logger.Error("userRepo.AddUser call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	// registration contant is by definition the prinary contact.
	// TODO: normalize contact principal
	newContact := models.NewContact(newUser.ID, "", contactPrincipal, contactType, true)
	err = us.contactRepo.AddContact(ctx, &newContact, initiator)
	if err != nil {
		logger.Error("contactRepo.AddContact call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	// generate confirmation code
	// TODO: make token valid time configurable
	confirmationToken, err := models.NewToken(newContact.ID, models.TokenTypeConfirmContact, time.Hour*2)
	if err != nil {
		evtString := "failed to create new contact confirmation token"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	err = us.tokenService.PutToken(ctx, logger, confirmationToken)
	if err != nil {
		evtString := "failed to store new contact confirmation token"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	// send confirmation email
	// TODO: convert this email into a template...
	to := []string{contactPrincipal}
	err = us.emailService.SendPlainTextEmail(ctx, logger, to, "contact confirmation link", confirmationToken.Value)
	if err != nil {
		evtString := "failed to send contact confirmation notification error occurred"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err // TODO: what should we do here???
	}
	// NOTE: allow user to set password on confirmation link click.
	span.AddEvent("user registered and confirmation notification sent")
	return nil
}

func (us userService) GetUserPrimaryContact(ctx context.Context, logger *zap.Logger, userID string, contactType string, initiator string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUserPrimaryContact")
	defer span.End()
	contact, err := us.contactRepo.GetPrimaryContactByUserID(ctx, userID, contactType)
	if err != nil {
		evtString := "failed to retreive primary contact by user id"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return models.Contact{}, err
	}
	span.AddEvent("user primary contact retreived")
	return contact, nil
}

func (us userService) GetUsersContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersContacts")
	defer span.End()
	contacts, err := us.contactRepo.GetContactsByUserID(ctx, userID)
	if err != nil {
		evtString := "failed to retreive contact by user id"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return nil, err
	}
	span.AddEvent("user contacts retreived")
	return contacts, nil
}

func (us userService) GetUsersConfirmedContacts(ctx context.Context, logger *zap.Logger, userID string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersConfirmedContacts")
	defer span.End()
	contacts, err := us.contactRepo.GetContactsByUserID(ctx, userID)
	if err != nil {
		evtString := "failed to retreive contact by user id"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return nil, err
	}
	confirmedContacts := make([]models.Contact, 0, len(contacts))
	for _, c := range contacts {
		if c.IsConfirmed() {
			confirmedContacts = append(confirmedContacts, c)
		}
	}
	// TODO: if there are no confirmed contact, should we return a no contacts found error or leave it as is?
	span.AddEvent("user confirmed contacts of type retreived")
	return confirmedContacts, nil
}

func (us userService) GetUsersContactsOfType(ctx context.Context, logger *zap.Logger, userID string, contactType string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersContactsOfType")
	defer span.End()
	contacts, err := us.contactRepo.GetContactsByUserIDAndType(ctx, userID, contactType)
	if err != nil {
		evtString := "failed to retreive contact by user id"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return nil, err
	}
	span.AddEvent("user contacts of type retreived")
	return contacts, nil
}

func (us userService) GetUsersConfirmedContactsOfType(ctx context.Context, logger *zap.Logger, userID string, contactType string, initiator string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "GetUsersConfirmedContactsOfType")
	defer span.End()
	contacts, err := us.contactRepo.GetContactsByUserIDAndType(ctx, userID, contactType)
	if err != nil {
		evtString := "failed to retreive contact by user id"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return nil, err
	}
	confirmedContacts := make([]models.Contact, 0, len(contacts))
	for _, c := range contacts {
		if c.IsConfirmed() {
			confirmedContacts = append(confirmedContacts, c)
		}
	}
	span.AddEvent("user confirmed contacts of type retreived")
	return confirmedContacts, nil
}

func (us userService) AddContact(ctx context.Context, logger *zap.Logger, userID string, contact *models.Contact, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "AddContact")
	defer span.End()
	// check user ids match
	if userID != contact.UserID {
		err := coreerrors.NewUserIDsDoNotMatchError(userID, contact.UserID, true)
		evtString := "user id and contact user id do not match"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// check that no other contact that is confirmed is already in data store
	err := us.checkForExistingConfirmedContacts(ctx, logger, &span, contact.Type, contact.Principal, userID)
	if err != nil {
		return err
	}

	// check if there is an existing primary contact of this type... if so throw an error.
	if contact.IsPrimary {
		_, err := us.contactRepo.GetPrimaryContactByUserID(ctx, userID, contact.Type)
		if err != nil && err.GetErrorCode() != coreerrors.ErrCodeNoContactFound {
			// we got an error that indicates there was an error
			// making the query against the data store
			evtString := "failed to retreive current primary contact of type"
			logger.Error(evtString, zap.Any("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return err
		}
		if err == nil {
			// there was no error so a contact was returned.
			// If you want to make a new primary contact you need to go through the
			// SetContactAsPrimary method
			err := coreerrors.NewContactToAddMarkedAsPrimaryError(userID, contact.Principal, contact.Type, true)
			evtString := "user id and contact user id do not match"
			logger.Error(evtString, zap.Any("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, evtString)
			return err
		}
		// the account has no contact marked as primary with the contact type
		logger.Warn("contact is marked as primary and there are no existing primary contacts of this type", zap.String("userID", userID), zap.String("contactPrincipal", contact.Principal), zap.String("contactType", contact.Type), zap.Bool("contactIsPrimary", contact.IsPrimary))
	}

	// run add contact on data store repo
	err = us.contactRepo.AddContact(ctx, contact, initiator)
	if err != nil {
		evtString := "failed to save contact in the data store"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	// TODO: send confirmation message?
	span.AddEvent("contact added for user")
	return nil
}

func (us userService) SetContactAsPrimary(ctx context.Context, logger *zap.Logger, userID string, newPrimaryContactID string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "SwapPrimaryContactOfType")
	defer span.End()
	span.AddEvent("primary contact swapped")
	return coreerrors.NewNotImplementedError(true)
}

func (us userService) ConfirmContact(ctx context.Context, logger *zap.Logger, confirmationCode string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "ConfirmContact")
	defer span.End()
	span.AddEvent("contact confirmed")
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

func (us userService) checkForExistingConfirmedContacts(ctx context.Context, logger *zap.Logger, span *trace.Span, contactType, contactPrincipal, userID string) errors.RichError {
	numExistingConfirmedContacts, err := us.contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType(ctx, contactType, contactPrincipal)
	if err != nil {
		logger.Error("contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(span, err, "")
		return err
	}
	if numExistingConfirmedContacts != 0 {
		if numExistingConfirmedContacts > 1 {
			// This is really bad and we need to know about it asap!
			errMsg := "critical issue here more than one contact is confirmed with this info"
			err = coreerrors.NewMultipleConfirmedInstancesOfContactError(contactPrincipal, contactType, numExistingConfirmedContacts, true)
			logger.Error(errMsg, zap.Any("error", err))
			apptelemetry.SetSpanOriginalError(span, err, errMsg)
		} else {
			errMsg := "a contact already exists and is confirmed with the data provided"
			errorFields := make(map[string]interface{})
			errorFields["numExistingConfirmedContacts"] = numExistingConfirmedContacts
			if userID == "" {
				err = coreerrors.NewRegistrationContactAlreadyConfirmedError(contactPrincipal, contactType, errorFields, true)
			} else {
				err = coreerrors.NewContactToAddAlreadyConfirmedError(userID, contactPrincipal, contactType, errorFields, true)
			}
			logger.Error(errMsg,
				zap.Any("error", err))
			apptelemetry.SetSpanOriginalError(span, err, errMsg)
		}
		return err
	}
	return nil
}
