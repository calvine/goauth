package service

import (
	"context"

	"github.com/calvine/goauth/core/apptelemetry"
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

func (appService) GetName() string {
	return "appService"
}

func (as appService) GetAppsByOwnerID(ctx context.Context, logger *zap.Logger, ownerID string, initiator string) ([]models.App, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetAppsByOwnerID")
	defer span.End()
	apps, err := as.appRepo.GetAppsByOwnerID(ctx, ownerID)
	if err != nil {
		logger.Error("appRepo.GetAppsByOwnerID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return nil, err
	}
	span.AddEvent("apps retreived")
	return apps, nil
}

func (as appService) GetAppByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.App, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetAppByID")
	defer span.End()
	app, err := as.appRepo.GetAppByID(ctx, id)
	if err != nil {
		logger.Error("appRepo.GetAppByID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.App{}, err
	}
	span.AddEvent("app retreived")
	return app, nil
}

func (as appService) GetAppByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetAppByClientID")
	defer span.End()
	app, err := as.appRepo.GetAppByClientID(ctx, clientID)
	if err != nil {
		logger.Error("appRepo.GetAppByClientID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.App{}, err
	}
	span.AddEvent("app retrevied")
	return app, nil
}

func (as appService) GetAppAndScopesByClientID(ctx context.Context, logger *zap.Logger, clientID string, initiator string) (models.App, []models.Scope, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetAppAndScopesByClientID")
	defer span.End()
	app, scopes, err := as.appRepo.GetAppAndScopesByClientID(ctx, clientID)
	if err != nil {
		logger.Error("appRepo.GetAppAndScopesByClientID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.App{}, nil, err
	}
	span.AddEvent("app and scopes retreived")
	return app, scopes, nil
}

func (as appService) AddApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "AddApp")
	defer span.End()
	err := models.ValidateApp(false, *app)
	if err != nil {
		evtString := "app data failed validation"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("app validated")
	err = as.appRepo.AddApp(ctx, app, initiator)
	if err != nil {
		logger.Error("appRepo.appRepo call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("app stored")
	return nil
}

func (as appService) UpdateApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "UpdateApp")
	defer span.End()
	err := models.ValidateApp(true, *app)
	if err != nil {
		evtString := "app data failed validation"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("app validated")
	err = as.appRepo.UpdateApp(ctx, app, initiator)
	if err != nil {
		logger.Error("appRepo.UpdateApp call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("app updated")
	return nil
}

func (as appService) DeleteApp(ctx context.Context, logger *zap.Logger, app *models.App, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "DeleteApp")
	defer span.End()
	err := as.appRepo.DeleteApp(ctx, app, initiator)
	if err != nil {
		logger.Error("appRepo.DeleteApp call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("app deleted")
	return nil
}

func (as appService) GetScopeByID(ctx context.Context, logger *zap.Logger, id string, initiator string) (models.Scope, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetScopeByID")
	defer span.End()
	scope, err := as.appRepo.GetScopeByID(ctx, id)
	if err != nil {
		logger.Error("appRepo.GetScopeByID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return models.Scope{}, err
	}
	span.AddEvent("scope retreived")
	return scope, nil
}

func (as appService) GetScopesByAppID(ctx context.Context, logger *zap.Logger, appID string, initiator string) ([]models.Scope, errors.RichError) {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "GetScopesByAppID")
	defer span.End()
	scopes, err := as.appRepo.GetScopesByAppID(ctx, appID)
	if err != nil {
		logger.Error("appRepo.GetScopesByAppID call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return nil, err
	}
	span.AddEvent("scopes retreived")
	return scopes, nil
}

// TODO: Determine if needed...
// func (as appService) GetScopesByClientID(ctx context.Context, clientID string, initiator string) ([]models.Scope, errors.RichError) {
// 	return nil, coreerrors.NewNotImplementedError(true)
// }

func (as appService) AddScopeToApp(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "AddScopeToApp")
	defer span.End()
	err := models.ValidateScope(false, *scope)
	if err != nil {
		evtString := "scope failed validation"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("scope validated")
	err = as.appRepo.AddScope(ctx, scope, initiator)
	if err != nil {
		logger.Error("appRepo.AddScope call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("scope stored")
	return nil
}

func (as appService) UpdateScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "UpdateScope")
	defer span.End()
	err := models.ValidateScope(false, *scope)
	if err != nil {
		evtString := "app scope failed validation"
		logger.Error(evtString, zap.Any("error", err))
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	span.AddEvent("scope validated")
	err = as.appRepo.UpdateScope(ctx, scope, initiator)
	if err != nil {
		logger.Error("appRepo.UpdateScope call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("scope updated")
	return nil
}

func (as appService) DeleteScope(ctx context.Context, logger *zap.Logger, scope *models.Scope, initiator string) errors.RichError {
	span := apptelemetry.CreateFunctionSpan(ctx, as.GetName(), "UpdateScope")
	defer span.End()
	err := as.appRepo.DeleteScope(ctx, scope, initiator)
	if err != nil {
		logger.Error("appRepo.DeleteScope call failed", zap.Any("error", err))
		apptelemetry.SetSpanError(&span, err, "")
		return err
	}
	span.AddEvent("scope deleted")
	return nil
}
