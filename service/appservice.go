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

func (as appService) GetAppsByOwnerID(ctx context.Context, id string, initiator string) ([]models.App, errors.RichError) {
	return nil, coreerrors.NewNotImplementedError(true)
}

func (as appService) GetAppByID(ctx context.Context, id string, initiator string) (models.App, errors.RichError) {
	return models.App{}, coreerrors.NewNotImplementedError(true)
}

func (as appService) GetAppByClientID(ctx context.Context, clientID string, initiator string) (models.App, errors.RichError) {
	return models.App{}, coreerrors.NewNotImplementedError(true)
}

func (as appService) GetAppAndScopesByClientID(ctx context.Context, clientID string, initiator string) (models.App, []models.Scope, errors.RichError) {
	return models.App{}, nil, nil
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
	return models.Scope{}, coreerrors.NewNotImplementedError(true)
}

func (as appService) GetScopesByAppID(ctx context.Context, appID string, initiator string) ([]models.Scope, errors.RichError) {
	return nil, coreerrors.NewNotImplementedError(true)
}

func (as appService) GetScopesByClientID(ctx context.Context, clientID string, initiator string) ([]models.Scope, errors.RichError) {
	return nil, coreerrors.NewNotImplementedError(true)
}

func (as appService) AddScopesToApp(ctx context.Context, scopes []*models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) UpdateScope(ctx context.Context, scope *models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}

func (as appService) DeleteScope(ctx context.Context, scope *models.Scope, initiator string) errors.RichError {
	return coreerrors.NewNotImplementedError(true)
}
