package repotest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
)

// TODO: improve repo tests to cover fail conditions, confirm data is manipulated properly in the underlying data store.

var (
	initialTestUser models.User

	initialTestContact models.Contact

	initialTestApp models.App

	initialTestApp2 models.App

	// initialTestAppClientSecret string

	// initialTestApp2ClientSecret string

	initialTestAppScopes []models.Scope

	// initialTestApp2Scopes []models.Scope
)

const (
	numScopes = 10
)

type RepoTestHarnessInput struct {
	UserRepo              *repo.UserRepo
	ContactRepo           *repo.ContactRepo
	AddressRepo           *repo.AddressRepo
	ProfileRepo           *repo.ProfileRepo
	AppRepo               *repo.AppRepo
	TokenRepo             *repo.TokenRepo
	AuditLogRepo          *repo.AuditLogRepo
	SetupTestDataSource   func(t *testing.T, input RepoTestHarnessInput)
	CleanupTestDataSource func(t *testing.T, input RepoTestHarnessInput)
}

func RunReposTestHarness(t *testing.T, implementationName string, input RepoTestHarnessInput) {

	if input.SetupTestDataSource != nil {
		input.SetupTestDataSource(t, input)
	}
	if input.CleanupTestDataSource != nil {
		defer input.CleanupTestDataSource(t, input)
	}
	setupTestHarnessData(t, input)

	// functionality tests
	t.Run(fmt.Sprintf("%s - userRepo", implementationName), func(t *testing.T) {
		if input.UserRepo != nil {
			testUserRepo(t, *input.UserRepo)
		} else {
			t.Skipf("no implementation for %s provided for userRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - contactRepo", implementationName), func(t *testing.T) {
		if input.ContactRepo != nil {
			testContactRepo(t, *input.ContactRepo)
		} else {
			t.Skipf("no implementation for %s provided for contactRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - addressRepo", implementationName), func(t *testing.T) {
		if input.AddressRepo != nil {
			testAddressRepo(t, *input.AddressRepo)
		} else {
			t.Skipf("no implementation for %s provided for addressRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - profileRepo", implementationName), func(t *testing.T) {
		if input.ProfileRepo != nil {
			testProfileRepo(t, *input.ProfileRepo)
		} else {
			t.Skipf("no implementation for %s provided for profileRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - appRepo", implementationName), func(t *testing.T) {
		if input.AppRepo != nil {
			testAppRepo(t, *input.AppRepo)
		} else {
			t.Skipf("no implementation for %s provided for appRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - tokenRepo", implementationName), func(t *testing.T) {
		if input.TokenRepo != nil {
			testTokenRepo(t, *input.TokenRepo)
		} else {
			t.Skipf("no implementation for %s provided for tokenRepo", implementationName)
		}
	})

	t.Run(fmt.Sprintf("%s - auditLogRepo", implementationName), func(t *testing.T) {
		if input.AuditLogRepo != nil {
			testAuditLogRepo(t, *input.AuditLogRepo)
		} else {
			t.Skipf("no implementation for %s provided for auditLogRepo", implementationName)
		}
	})
}

func setupTestHarnessData(t *testing.T, input RepoTestHarnessInput) {
	var err errors.RichError
	createdByID := "test setup"
	// create test user
	initialTestUser = models.NewUser()
	initialTestUser.PasswordHash = "passwordhash2"
	initialTestUser.LastLoginDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	// add a test user
	err = (*input.UserRepo).AddUser(context.TODO(), &initialTestUser, createdByID)
	if err != nil {
		t.Errorf("setup failed to add user to database: %s", err.Error())
	}

	// create test contact
	initialTestContact = models.NewContact(initialTestUser.ID, "", "InitialTestUser@email.com", core.CONTACT_TYPE_EMAIL, true)
	// add a test contact for the test user.
	err = (*input.ContactRepo).AddContact(context.TODO(), &initialTestContact, createdByID)
	if err != nil {
		t.Errorf("setup failed to add contact to database: %s", err.Error())
	}

	if input.AppRepo != nil {
		// create test apps
		initialTestApp, _, err = models.NewApp(initialTestUser.ID, "test app 1", "https://my.app/callback", "https://my.app/assets/logo.png")
		if err != nil {
			t.Errorf("setup failed to create app: %s", err.Error())
		}
		initialTestApp2, _, err = models.NewApp(initialTestUser.ID, "test app 1", "https://my.app/callback", "https://my.app/assets/logo.png")
		if err != nil {
			t.Errorf("setup failed to create app: %s", err.Error())
		}
		// add test apps
		err = (*input.AppRepo).AddApp(context.TODO(), &initialTestApp, createdByID)
		if err != nil {
			t.Errorf("setup failed to add app to database: %s", err.Error())
		}
		err = (*input.AppRepo).AddApp(context.TODO(), &initialTestApp2, createdByID)
		if err != nil {
			t.Errorf("setup failed to add app to database: %s", err.Error())
		}
		// create test scopes
		initialTestAppScopes = make([]models.Scope, 0, numScopes)
		// initialTestApp2Scopes = make([]models.Scope, 0, numScopes)
		// add scopes to test app
		for i := 1; i <= numScopes; i++ {
			scope := models.NewScope(initialTestApp.ID, fmt.Sprintf("app_scope_%d", i), fmt.Sprintf("permissions associated with app_scope_%d", i))
			err = (*input.AppRepo).AddScope(context.TODO(), &scope, createdByID)
			if err != nil {
				t.Errorf("setup failed to add scope %d to database: %s", i, err.Error())
			}
			initialTestAppScopes = append(initialTestAppScopes, scope)

			// scopes will be added to initialTestApp2 in the add scope test on the app repo tests.

			// scope2 := models.NewScope(initialTestApp2.ID, fmt.Sprintf("other_app_scope_%d", i), fmt.Sprintf("permissions associated with other_app_scope_%d", i))
			// err = (*input.AppRepo).AddScope(context.TODO(), &scope2, createdByID)
			// if err != nil {
			// 	t.Errorf("setup failed to add other app scope %d to database: %s", i, err.Error())
			// }
			// initialTestApp2Scopes = append(initialTestApp2Scopes, scope2)
		}
	}
}
