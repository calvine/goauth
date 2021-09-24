package memory

import (
	"context"
	"time"

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
	app, ok := (*ar.apps)[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		return app, coreerrors.NewNoAppFoundError(fields, true)
	}
	return app, nil
}

func (ar appRepo) GetAppByClientID(ctx context.Context, clientID string) (models.App, errors.RichError) {
	var app models.App
	for _, a := range *ar.apps {
		if a.ClientID == clientID {
			app = a
			break
		}
	}
	return app, nil
}

func (ar appRepo) GetAppsByOwnerID(ctx context.Context, ownerID string) ([]models.App, errors.RichError) {
	apps := make([]models.App, 0)
	for _, app := range *ar.apps {
		if app.OwnerID == ownerID {
			apps = append(apps, app)
		}
	}
	if len(apps) == 0 {
		fields := map[string]interface{}{"ownerID": ownerID}
		return apps, coreerrors.NewNoAppFoundError(fields, true)
	}
	return apps, nil
}

func (ar appRepo) GetAppAndScopesByClientID(ctx context.Context, clientID string) (models.App, []models.Scope, errors.RichError) {
	var app models.App
	var scopes []models.Scope
	for _, a := range *ar.apps {
		if a.ClientID == clientID {
			app = a
			scopes = (*ar.appScopes)[a.ID]
			return app, scopes, nil
		}
	}
	fields := map[string]interface{}{"clientID": clientID}
	return app, scopes, coreerrors.NewNoAppFoundError(fields, true)
}

func (ar appRepo) AddApp(ctx context.Context, app *models.App, createdBy string) errors.RichError {
	app.AuditData.CreatedByID = createdBy
	app.AuditData.CreatedOnDate = time.Now().UTC()
	if app.ID == "" {
		app.ID = uuid.Must(uuid.NewRandom()).String()
	}
	(*ar.apps)[app.ID] = *app
	(*ar.appScopes)[app.ID] = make([]models.Scope, 0, 5)
	return nil
}

func (ar appRepo) UpdateApp(ctx context.Context, app *models.App, modifiedBy string) errors.RichError {
	app.AuditData.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedBy}
	app.AuditData.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	(*ar.apps)[app.ID] = *app
	return nil
}

func (ar appRepo) DeleteApp(ctx context.Context, app *models.App, deletedBy string) errors.RichError {
	_, ok := (*ar.apps)[app.ID]
	if !ok {
		fields := map[string]interface{}{"id": app.ID}
		return coreerrors.NewNoAppFoundError(fields, true)
	}
	delete(*ar.apps, app.ID)
	delete(*ar.appScopes, app.ID)
	return nil
}

func (ar appRepo) GetScopeByID(ctx context.Context, id string) (models.Scope, errors.RichError) {
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
		return models.Scope{}, coreerrors.NewNoScopeFoundError(fields, true)
	}
	return scope, nil
}

func (ar appRepo) GetScopesByAppID(ctx context.Context, appID string) ([]models.Scope, errors.RichError) {
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return scopes, coreerrors.NewNoScopeFoundError(fields, true)
	}
	return scopes, nil
}

func (ar appRepo) AddScope(ctx context.Context, scope *models.Scope, createdBy string) errors.RichError {
	appID := scope.AppID
	scope.AuditData.CreatedByID = createdBy
	scope.AuditData.CreatedOnDate = time.Now().UTC()
	if scope.ID == "" {
		scope.ID = uuid.Must(uuid.NewRandom()).String()
	}
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"ID": scope.ID, "AppID": scope.AppID}
		return coreerrors.NewNoAppFoundError(fields, true)
	}
	scopes = append(scopes, *scope)
	(*ar.appScopes)[appID] = scopes
	return nil
}

func (ar appRepo) UpdateScope(ctx context.Context, scope *models.Scope, modifiedBy string) errors.RichError {
	appID := scope.AppID
	scopeID := scope.ID
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return coreerrors.NewNoScopeFoundError(fields, true)
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
		return coreerrors.NewNoScopeFoundError(fields, true)
	}
	return nil
}

func (ar appRepo) DeleteScope(ctx context.Context, scope *models.Scope, deletedBy string) errors.RichError {
	appID := scope.AppID
	scopeID := scope.ID
	scopes, ok := (*ar.appScopes)[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return coreerrors.NewNoScopeFoundError(fields, true)
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
		return coreerrors.NewNoScopeFoundError(fields, true)
	}
	return nil
}
