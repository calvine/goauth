package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

type appRepo struct {
	apps      *map[string]models.App
	appScopes *map[string][]models.Scope
}

func NewMemoryAppRepo() repo.AppRepo {
	apps = make(map[string]models.App)
	appScopes = make(map[string][]models.Scope)
	return appRepo{
		apps:      &apps,
		appScopes: &appScopes,
	}
}

func (appRepo) GetName() string {
	return "appRepo"
}

func (appRepo) GetType() string {
	return dataSourceType
}

func (ar appRepo) GetAppByID(ctx context.Context, id string) (models.App, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetAppByID", ar.GetType())
	defer span.End()
	app, ok := (*ar.apps)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		err := coreerrors.NewNoAppFoundError(fields, true)
		evtString := fmt.Sprintf("no app found with id: %s", id)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return app, err
	}
	span.AddEvent("app retreived")
	return app, nil
}

func (ar appRepo) GetAppByClientID(ctx context.Context, clientID string) (models.App, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetAppByClientID", ar.GetType())
	defer span.End()
	var app models.App
	for _, a := range *ar.apps {
		if a.ClientID == clientID {
			app = a
			break
		}
	}
	span.AddEvent("app retreived")
	return app, nil
}

func (ar appRepo) GetAppsByOwnerID(ctx context.Context, ownerID string) ([]models.App, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetAppsByOwnerID", ar.GetType())
	defer span.End()
	apps := make([]models.App, 0)
	for _, app := range *ar.apps {
		if app.OwnerID == ownerID {
			apps = append(apps, app)
		}
	}
	if len(apps) == 0 {
		fields := map[string]interface{}{"ownerID": ownerID}
		err := coreerrors.NewNoAppFoundError(fields, true)
		evtString := fmt.Sprintf("no apps found for owner id: %s", ownerID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return apps, err
	}
	span.AddEvent("apps retreived")
	return apps, nil
}

func (ar appRepo) GetAppAndScopesByClientID(ctx context.Context, clientID string) (models.App, []models.Scope, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetAppAndScopesByClientID", ar.GetType())
	defer span.End()
	var app models.App
	var scopes []models.Scope
	for _, a := range *ar.apps {
		if a.ClientID == clientID {
			app = a
			scopes = (*ar.appScopes)[a.ID]
			span.AddEvent("app and scopes retreived")
			return app, scopes, nil
		}
	}
	fields := map[string]interface{}{"clientID": clientID}
	err := coreerrors.NewNoAppFoundError(fields, true)
	evtString := fmt.Sprintf("no apps or scopes found for client id: %s", clientID)
	span.AddEvent(evtString)
	apptelemetry.SetSpanError(&span, err)
	return app, scopes, err
}

func (ar appRepo) AddApp(ctx context.Context, app *models.App, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "AddApp", ar.GetType())
	defer span.End()
	app.AuditData.CreatedByID = createdBy
	app.AuditData.CreatedOnDate = time.Now().UTC()
	if app.ID == "" {
		app.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*ar.apps)[app.ID] = *app
	(*ar.appScopes)[app.ID] = make([]models.Scope, 0, 5)
	span.AddEvent("app stored")
	return nil
}

func (ar appRepo) UpdateApp(ctx context.Context, app *models.App, modifiedBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "UpdateApp", ar.GetType())
	defer span.End()
	app.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedBy}
	app.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*ar.apps)[app.ID] = *app
	span.AddEvent("app updated")
	return nil
}

func (ar appRepo) DeleteApp(ctx context.Context, app *models.App, deletedBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "DeleteApp", ar.GetType())
	defer span.End()
	_, ok := (*ar.apps)[app.ID]
	if !ok {
		fields := map[string]interface{}{"id": app.ID}
		err := coreerrors.NewNoAppFoundError(fields, true)
		evtString := fmt.Sprintf("no app found with id: %s", app.ID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	delete(*ar.apps, app.ID)
	delete(*ar.appScopes, app.ID)
	span.AddEvent("app and scopes deleted")
	return nil
}

func (ar appRepo) GetScopeByID(ctx context.Context, id string) (models.Scope, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetScopeByID", ar.GetType())
	defer span.End()
	var scope models.Scope
	found := false
	for _, appScopes := range *ar.appScopes {
		for _, s := range appScopes {
			if s.ID == id {
				scope = s
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		fields := map[string]interface{}{"id": id}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scope found with id: %s", id)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return models.Scope{}, err
	}
	span.AddEvent("scope retreived")
	return scope, nil
}

func (ar appRepo) GetScopesByAppID(ctx context.Context, appID string) ([]models.Scope, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "GetScopesByAppID", ar.GetType())
	defer span.End()
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scopes found with app id: %s", appID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return scopes, err
	}
	span.AddEvent("scopes retrevied")
	return scopes, nil
}

func (ar appRepo) AddScope(ctx context.Context, scope *models.Scope, createdBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "AddScope", ar.GetType())
	defer span.End()
	appID := scope.AppID
	scope.AuditData.CreatedByID = createdBy
	scope.AuditData.CreatedOnDate = time.Now().UTC()
	if scope.ID == "" {
		scope.ID = uuid.Must(uuid.NewRandom()).String()
	}
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"ID": scope.ID, "AppID": scope.AppID}
		err := coreerrors.NewNoAppFoundError(fields, true)
		evtString := fmt.Sprintf("no app found with app id: %s", scope.AppID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	scopes = append(scopes, *scope)
	(*ar.appScopes)[appID] = scopes
	span.AddEvent("scope stored")
	return nil
}

func (ar appRepo) UpdateScope(ctx context.Context, scope *models.Scope, modifiedBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "UpdateScope", ar.GetType())
	defer span.End()
	appID := scope.AppID
	scopeID := scope.ID
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scope found with id: %s", scope.ID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	scopeFound := false
	for i, s := range scopes {
		if s.ID == scopeID {
			scope.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedBy}
			scope.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
			scopes[i] = *scope
			(*ar.appScopes)[appID] = scopes
			scopeFound = true
			break
		}
	}
	if !scopeFound {
		fields := map[string]interface{}{"ID": scope.ID, "appID": scope.AppID}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scope found with id: %s", scope.ID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	span.AddEvent("scope updated")
	return nil
}

func (ar appRepo) DeleteScope(ctx context.Context, scope *models.Scope, deletedBy string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ar.GetName(), "DeleteScope", ar.GetType())
	defer span.End()
	appID := scope.AppID
	scopeID := scope.ID
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scope found with id: %s", scope.ID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	scopeFound := false
	for i, s := range scopes {
		if s.ID == scopeID {
			// remove specific item from slice DOES NOT PRESERVE ORDER....
			scopes[i] = scopes[len(scopes)-1]
			scopes = scopes[:len(scopes)-1]
			// same, but would preserve order
			// scopes = append(scopes[:i], scopes[i+1:]...)
			scope = nil
			(*ar.appScopes)[appID] = scopes
			scopeFound = true
			break
		}
	}
	if !scopeFound {
		fields := map[string]interface{}{"appID": appID}
		err := coreerrors.NewNoScopeFoundError(fields, true)
		evtString := fmt.Sprintf("no scope found with id: %s", scope.ID)
		span.AddEvent(evtString)
		apptelemetry.SetSpanError(&span, err)
		return err
	}
	span.AddEvent("scope deleted")
	return nil
}
