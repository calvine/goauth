package service

// type userService struct{}

// func NewUserService() services.UserService {

// }

// func (us userService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, errors.RichError) {
// 	token, err := ls.tokenService.GetToken(confirmationCode, models.TokenTypeConfirmContact)
// 	if err != nil {
// 		return false, err
// 	}
// 	contact, err := ls.contactRepo.GetContactByContactId(ctx, token.TargetID)
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
