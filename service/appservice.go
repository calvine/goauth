package service

import (
	"context"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/services"
	"github.com/calvine/richerror/errors"
)

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

func (as appService) GetAppsByOwnerID(ctx context.Context, ownerID string, initiator string) ([]models.App, errors.RichError) {
	apps, err := as.appRepo.GetAppsByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (as appService) GetAppByID(ctx context.Context, id string, initiator string) (models.App, errors.RichError) {
	app, err := as.appRepo.GetAppByID(ctx, id)
	if err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (as appService) GetAppByClientID(ctx context.Context, clientID string, initiator string) (models.App, errors.RichError) {
	app, err := as.appRepo.GetAppByClientID(ctx, clientID)
	if err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (as appService) GetAppAndScopesByClientID(ctx context.Context, clientID string, initiator string) (models.App, []models.Scope, errors.RichError) {
	app, scopes, err := as.appRepo.GetAppAndScopesByClientID(ctx, clientID)
	if err != nil {
		return models.App{}, nil, err
	}
	return app, scopes, nil
}

func (as appService) AddApp(ctx context.Context, app *models.App, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) UpdateApp(ctx context.Context, app *models.App, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) DeleteApp(ctx context.Context, app *models.App, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) GetScopeByID(ctx context.Context, id string, initiator string) (models.Scope, errors.RichError) {
	scope, err := as.appRepo.GetScopeByID(ctx, id)
	if err != nil {
		return models.Scope{}, err
	}
	return scope, nil
}

func (as appService) GetScopesByAppID(ctx context.Context, appID string, initiator string) ([]models.Scope, errors.RichError) {
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

func (as appService) AddScopeToApp(ctx context.Context, scopes *models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) UpdateScope(ctx context.Context, scope *models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) DeleteScope(ctx context.Context, scope *models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}
