package memory

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

type appRepo struct {
	apps      map[string](models.App)
	appScopes map[string][]models.Scope
}

func NewMemoryAppRepo() repo.AppRepo {
	apps := make(map[string]models.App)
	appScopes := make(map[string][]models.Scope)
	return appRepo{
		apps:      apps,
		appScopes: appScopes,
	}
}

func (ar appRepo) GetAppByID(ctx context.Context, id string) (models.App, errors.RichError) {
	app, ok := ar.apps[id]
	if !ok {
		fields := map[string]interface{}{"id": id}
		return app, coreerrors.NewNoAppFoundError(fields, true)
	}
	return app, nil
}
func (ar appRepo) GetAppsByOwnerID(ctx context.Context, ownerID string) ([]models.App, errors.RichError) {
	apps := make([]models.App, 0)
	for _, app := range ar.apps {
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
func (ar appRepo) GetAppAndScopesByClientIDAndCallbackURI(ctx context.Context, clientID, callbackURI string) (models.App, []models.Scope, errors.RichError) {
	var app models.App
	var scopes []models.Scope
	for _, a := range ar.apps {
		if a.ClientID == clientID && a.CallbackURI == callbackURI {
			app = a
			scopes = ar.appScopes[a.ID]
			return app, scopes, nil
		}
	}
	fields := map[string]interface{}{"clientID": clientID, "callbackURI": callbackURI}
	return app, scopes, coreerrors.NewNoAppFoundError(fields, true)
}
func (ar appRepo) AddApp(ctx context.Context, app *models.App, createdBy string) errors.RichError {
	app.CreatedByID = createdBy
	app.CreatedOnDate = time.Now().UTC()
	ar.apps[app.ID] = *app
	return nil
}
func (ar appRepo) UpdateApp(ctx context.Context, app *models.App, modifiedBy string) errors.RichError {
	app.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedBy}
	app.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	ar.apps[app.ID] = *app
	return nil
}
func (ar appRepo) DeleteApp(ctx context.Context, app *models.App, deletedBy string) errors.RichError {
	delete(ar.apps, app.ID)
	delete(ar.appScopes, app.ID)
	return nil
}

func (ar appRepo) GetScopesByAppID(ctx context.Context, appID string) ([]models.Scope, errors.RichError) {
	scopes, ok := ar.appScopes[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return scopes, coreerrors.NewNoScopeFoundError(fields, true)
	}
	return scopes, nil
}
func (ar appRepo) AddScope(ctx context.Context, scope *models.Scope, createdBy string) errors.RichError {
	appID := scope.ApplicationID
	scope.CreatedByID = createdBy
	scope.CreatedOnDate = time.Now().UTC()
	scopes, ok := ar.appScopes[appID]
	if !ok {
		scopes = make([]models.Scope, 0, 1)
	}
	scopes = append(scopes, *scope)
	ar.appScopes[appID] = scopes
	return nil
}
func (ar appRepo) UpdateScope(ctx context.Context, scope *models.Scope, modifiedBy string) errors.RichError {
	appID := scope.ApplicationID
	scopeID := scope.ID
	scopes, ok := ar.appScopes[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return coreerrors.NewNoScopeFoundError(fields, true)
	}
	for _, s := range scopes {
		if s.ID == scopeID {
			scope.ModifiedByID = nullable.NullableString{HasValue: true, Value: modifiedBy}
			scope.ModifiedOnDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
			scopes = append(scopes, *scope)
			ar.appScopes[appID] = scopes
		}
	}
	return nil
}
func (ar appRepo) DeleteScope(ctx context.Context, scope *models.Scope, deletedBy string) errors.RichError {
	appID := scope.ApplicationID
	scopeID := scope.ID
	scopes, ok := ar.appScopes[appID]
	if !ok {
		fields := map[string]interface{}{"appID": appID}
		return coreerrors.NewNoScopeFoundError(fields, true)
	}
	for i, s := range scopes {
		if s.ID == scopeID {
			// remove specific item from slice DOES NOT PRESERVE ORDER....
			scopes[i] = scopes[len(scopes)-1]
			scopes = scopes[:len(scopes)-1]
			// same, but would preserve order
			scopes = append(scopes[:i], scopes[i+1:]...)
		}
	}
	return nil
}
