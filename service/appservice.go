package service

import (
	"context"

	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
	"go.uber.org/zap"
)

// TODO: need validated info added to mutative function parameters for another level of checking

type appService struct {
	appRepo      repo.AppRepo
	auditLogRepo repo.AuditLogRepo
}

func NewAppService(appRepo repo.AppRepo, auditLogRepo repo.AuditLogRepo) services.AppService {
	return appService{
		appRepo:      appRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (as appService) GetAppsByOwnerID(ctx context.Context, logger *zap.Logger, ownerID string, initiator string) ([]models.App, errors.RichError) {
	apps, err := as.appRepo.GetAppsByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (as appService) GetAppByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.App, errors.RichError) {
	app, err := as.appRepo.GetAppByID(ctx, id)
	if err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (as appService) GetAppByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, errors.RichError) {
	app, err := as.appRepo.GetAppByClientID(ctx, clientID)
	if err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (as appService) GetAppAndScopesByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, []models.Scope, errors.RichError) {
	app, scopes, err := as.appRepo.GetAppAndScopesByClientID(ctx, clientID)
	if err != nil {
		return models.App{}, nil, err
	}
	return app, scopes, nil
}

func (as appService) AddApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	err := models.ValidateApp(false, *app)
	if err != nil {
		return err
	}
	err = as.appRepo.AddApp(ctx, app, initiator)
	return err
}

func (as appService) UpdateApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	err := models.ValidateApp(true, *app)
	if err != nil {
		return err
	}
	err = as.appRepo.UpdateApp(ctx, app, initiator)
	return err
}

func (as appService) DeleteApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	err := as.appRepo.DeleteApp(ctx, app, initiator)
	return err
}

func (as appService) GetScopeByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.Scope, errors.RichError) {
	scope, err := as.appRepo.GetScopeByID(ctx, id)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
}

func (as appService) GetScopesByAppID(ctx context.Context, logger *zap.Logger, appID string, initiator string) ([]models.Scope, errors.RichError) {
	scopes, err := as.appRepo.GetScopesByAppID(ctx, appID)
	if err != nil {
		return nil, err
	}
	return scopes, nil
}

// TODO: Determine if needed...
// func (as appService) GetScopesByClientID(ctx context.Context, clientID string, initiator string) ([]models.Scope, errors.RichError) {
// 	return nil, coreerrors.NewNotImplementedError(true)
// }

func (as appService) AddScopeToApp(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	err := models.ValidateScope(false, *scope)
	if err != nil {
		return err
	}
	err = as.appRepo.AddScope(ctx, scope, initiator)
	return err
}

func (as appService) UpdateScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	err := models.ValidateScope(false, *scope)
	if err != nil {
		return err
	}
	err = as.appRepo.UpdateScope(ctx, scope, initiator)
	return err
}

func (as appService) DeleteScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	err := as.appRepo.DeleteScope(ctx, scope, initiator)
	return err
}
