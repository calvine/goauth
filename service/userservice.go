package service

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/models/email"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/goauth/internal/constants"
	"github.com/calvine/richerror/errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	contactConfirmationDuration = time.Hour * 4
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
		logger.Error("userRepo.GetUserAndContactByContact call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.User{}, models.Contact{}, err
	}
	if !contact.IsConfirmed() {
		evtString := fmt.Sprintf("contact found is not confirmed: ID = %s", contact.ID)
		fields := make(map[string]interface{})
		fields["userId"] = contact.UserID
		err := coreerrors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, fields, true)
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return models.User{}, models.Contact{}, err
	}
	span.AddEvent("user and contact retreived")
	return user, contact, nil
}

func (us userService) RegisterUserAndPrimaryContact(ctx context.Context, logger *zap.Logger, contactType, contactPrincipal string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "RegisterUserAndPrimaryContact")
	defer span.End()

	err := models.IsValidContactType(contactType)
	if err != nil {
		evtString := "failed to create new contact due to invalid contact type"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// registration contact is by definition the primary contact.
	newContact := models.NewContact("", "", contactPrincipal, contactType, true)
	err = models.IsValidNormalizedContactPrincipal(newContact.Type, newContact.Principal)
	if err != nil {
		evtString := "failed to create new contact due to invalid contact principal"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// TODO: normalize contact principal
	// check that email address does not already exist as a confirmed contact.
	err = us.checkForExistingConfirmedContacts(ctx, logger, &span, contactType, contactPrincipal, "")
	if err != nil {
		// additional error stuff handeled in checkForExistingConfirmedContacts function
		return err
	}
	// create new user and contact in datastore
	newUser := models.NewUser()
	err = us.userRepo.AddUser(ctx, &newUser, initiator)
	if err != nil {
		logger.Error("userRepo.AddUser call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}

	newContact.UserID = newUser.ID
	err = us.contactRepo.AddContact(ctx, &newContact, initiator)
	if err != nil {
		logger.Error("contactRepo.AddContact call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	// generate confirmation code
	// TODO: make token valid time configurable
	confirmationToken, err := models.NewToken(newContact.ID, models.TokenTypeConfirmContact, contactConfirmationDuration)
	if err != nil {
		evtString := "failed to create new contact confirmation token"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	err = us.tokenService.PutToken(ctx, logger, confirmationToken)
	if err != nil {
		evtString := "failed to store new contact confirmation token"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	// send confirmation email
	// TODO: convert this email into a template...
	emailMessage := email.EmailMessage{
		From:    constants.NoReplyEmailAddress,
		To:      []string{contactPrincipal},
		Subject: "contact confirmation link",
		Body:    confirmationToken.Value,
	}
	err = us.emailService.SendPlainTextEmail(ctx, logger, emailMessage)
	if err != nil {
		evtString := "failed to send contact confirmation notification error occurred"
		logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
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
	err := models.IsValidContactType(contact.Type)
	if err != nil {
		evtString := "failed to create new contact due to invalid contact type"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	err = models.IsValidNormalizedContactPrincipal(contact.Type, contact.Principal)
	if err != nil {
		evtString := "failed to create new contact due to invalid contact principal"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// check user ids match
	if userID != contact.UserID {
		err := coreerrors.NewUserIDsDoNotMatchError(userID, contact.UserID, true)
		evtString := "user id and contact user id do not match"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// check that no other contact that is confirmed is already in data store
	err = us.checkForExistingConfirmedContacts(ctx, logger, &span, contact.Type, contact.Principal, userID)
	if err != nil {
		// additional error stuff handeled in checkForExistingConfirmedContacts function
		return err
	}

	// check if there is an existing primary contact of this type... if so throw an error.
	if contact.IsPrimary {
		_, err := us.contactRepo.GetPrimaryContactByUserID(ctx, userID, contact.Type)
		if err != nil && err.GetErrorCode() != coreerrors.ErrCodeNoContactFound {
			// we got an error that indicates there was an error
			// making the query against the data store
			evtString := "failed to retreive current primary contact of type"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return err
		}
		if err == nil {
			// there was no error so a contact was returned.
			// If you want to make a new primary contact you need to go through the
			// SetContactAsPrimary method
			err := coreerrors.NewContactToAddMarkedAsPrimaryError(userID, contact.Principal, contact.Type, true)
			evtString := "user id and contact user id do not match"
			logger.Error(evtString, zap.Reflect("error", err))
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
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	// generate contact confirmation token
	// TODO: make token valid time configurable
	confirmationToken, err := models.NewToken(contact.ID, models.TokenTypeConfirmContact, contactConfirmationDuration)
	if err != nil {
		evtString := "failed to create new contact confirmation token"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	err = us.tokenService.PutToken(ctx, logger, confirmationToken)
	if err != nil {
		evtString := "failed to store new contact confirmation token"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	// send confirmation email
	// TODO: convert this email into a template...
	emailMessage := email.EmailMessage{
		From:    constants.NoReplyEmailAddress,
		To:      []string{contact.Principal},
		Subject: "contact confirmation link",
		Body:    confirmationToken.Value,
	}
	err = us.emailService.SendPlainTextEmail(ctx, logger, emailMessage)
	if err != nil {
		evtString := "failed to send contact confirmation notification error occurred"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err // TODO: what should we do here???
	}
	span.AddEvent("contact added for user")
	return nil
}

func (us userService) SetContactAsPrimary(ctx context.Context, logger *zap.Logger, userID string, newPrimaryContactID string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "SwapPrimaryContactOfType")
	defer span.End()
	newPrimaryContact, err := us.contactRepo.GetContactByID(ctx, newPrimaryContactID)
	if err != nil {
		evtString := "failed to retreive new primary contact of type"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	if newPrimaryContact.UserID != userID {
		err := coreerrors.NewUserIDsDoNotMatchError(userID, newPrimaryContact.UserID, true)
		evtString := "user id provided does not match user id of contact to set as primary"
		logger.Error(evtString, zap.String("userId", userID), zap.String("newPrimaryContactUserId", newPrimaryContact.UserID))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	if !newPrimaryContact.IsConfirmed() {
		fields := make(map[string]interface{})
		fields["userId"] = userID
		fields["isPrimary"] = newPrimaryContact.IsPrimary
		err := coreerrors.NewContactNotConfirmedError(newPrimaryContact.ID, newPrimaryContact.Principal, newPrimaryContact.Type, fields, true)
		evtString := "user id provided does not match user id of contact to set as primary"
		logger.Error(evtString, zap.String("userId", userID), zap.String("newPrimaryContactUserId", newPrimaryContact.UserID))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	hasCurrentPrimaryContact := true
	currentPrimaryContact, err := us.contactRepo.GetPrimaryContactByUserID(ctx, userID, newPrimaryContact.Type)
	if err != nil {
		if err.GetErrorCode() != coreerrors.ErrCodeNoContactFound {
			hasCurrentPrimaryContact = false
		} else {
			evtString := "failed to retreive current primary contact of type"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return err
		}
	}
	if hasCurrentPrimaryContact {
		if currentPrimaryContact.ID == newPrimaryContact.ID {
			err := coreerrors.NewContactAlreadyMarkedPrimaryError(currentPrimaryContact.Principal, currentPrimaryContact.Type, true)
			evtString := "contact to mark is primary is already primary contact"
			logger.Error(evtString, zap.String("userId", userID), zap.String("newPrimaryContactUserId", newPrimaryContact.UserID), zap.String("currentPrimaryContactUserId", currentPrimaryContact.UserID), zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(&span, err, evtString)
			return err
		}
		// we need to set the current contact to be not primary
		err = us.contactRepo.SwapPrimaryContacts(ctx, &currentPrimaryContact, &newPrimaryContact, initiator)
		if err != nil {
			evtString := "failed to swap primary state of the two contacts"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return err
		}
	} else {
		logger.Info("no current primary contact for type, updating contact provided as primary")
		newPrimaryContact.IsPrimary = true
		err := us.contactRepo.UpdateContact(ctx, &newPrimaryContact, initiator)
		if err != nil {
			evtString := "failed to update contact to set is primary flag"
			logger.Error(evtString, zap.Reflect("error", err))
			apptelemetry.SetSpanError(&span, err, evtString)
			return err
		}
	}

	span.AddEvent("primary contact swapped")
	return nil
}

func (us userService) ConfirmContact(ctx context.Context, logger *zap.Logger, confirmationCode string, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, us.GetName(), "ConfirmContact")
	defer span.End()
	confirmationToken, err := us.tokenService.GetToken(ctx, logger, confirmationCode, models.TokenTypeConfirmContact)
	if err != nil {
		evtString := "failed to retreive confirmation token from data store"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	if confirmationToken.TokenType != models.TokenTypeConfirmContact {
		err := coreerrors.NewInvalidTokenError(confirmationToken.Value, true)
		evtString := "token type is not valid"
		logger.Error(evtString, zap.String("tokenType", confirmationToken.TokenType.String()), zap.String("tokenValue", confirmationToken.Value), zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	// it appears that the token service will return this error if the token is expired, so this code is redundant...
	// should this service rely on the token service for some business logic?
	// I need to think on this...
	// if confirmationToken.IsExpired() {
	// 	err := coreerrors.NewExpiredTokenError(confirmationToken.Value, confirmationToken.TokenType.String(), confirmationToken.Expiration, true)
	// 	evtString := "token is expired"
	// 	logger.Error(evtString, zap.String("tokenType", confirmationToken.TokenType.String()), zap.String("tokenValue", confirmationToken.Value), zap.Reflect("error", err))
	// 	apptelemetry.SetSpanOriginalError(&span, err, evtString)
	// 	return err
	// }
	contactToConfirm, err := us.contactRepo.GetContactByID(ctx, confirmationToken.TargetID)
	if err != nil {
		evtString := "failed to retreive contact to confirm from data store"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	if contactToConfirm.IsConfirmed() {
		err := coreerrors.NewContactAlreadyConfirmedError(contactToConfirm.UserID, contactToConfirm.ID, contactToConfirm.Principal, contactToConfirm.Type, nil, true)
		evtString := "token is expired"
		logger.Error(evtString, zap.String("tokenType", confirmationToken.TokenType.String()), zap.String("tokenValue", confirmationToken.Value), zap.Reflect("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	contactToConfirm.ConfirmedDate.Set(time.Now().UTC())
	err = us.contactRepo.UpdateContact(ctx, &contactToConfirm, initiator)
	if err != nil {
		evtString := "failed to update contact to confirmed"
		logger.Error(evtString, zap.Reflect("error", err))
		apptelemetry.SetSpanError(&span, err, evtString)
		return err
	}
	span.AddEvent("contact confirmed")
	return nil
}

func (us userService) checkForExistingConfirmedContacts(ctx context.Context, logger *zap.Logger, span *trace.Span, contactType, contactPrincipal, userID string) errors.RichError {
	numExistingConfirmedContacts, err := us.contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType(ctx, contactType, contactPrincipal)
	if err != nil {
		logger.Error("contactRepo.GetExistingConfirmedContactsCountByPrincipalAndType call failed", zap.Reflect("error", err))
		apptelemetry.SetSpanError(span, err, "")
		return err
	}
	if numExistingConfirmedContacts != 0 {
		if numExistingConfirmedContacts > 1 {
			// This is really bad and we need to know about it asap!
			errMsg := "critical issue here more than one contact is confirmed with this info"
			err = coreerrors.NewMultipleConfirmedInstancesOfContactError(contactPrincipal, contactType, numExistingConfirmedContacts, true)
			logger.Error(errMsg, zap.Reflect("error", err))
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
			logger.Error(errMsg, zap.Reflect("error", err))
			apptelemetry.SetSpanOriginalError(span, err, errMsg)
		}
		return err
	}
	return nil
}
