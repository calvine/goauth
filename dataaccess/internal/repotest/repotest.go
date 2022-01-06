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

// TODO: improve repo tests to cover fail conditions, confirm data is manipulated properly in the underlying data store. Migrate to table tests

var (
	initialTestUser models.User

	initialTestUnconfirmedContact models.Contact

	initialTestConfirmedPrimaryContact models.Contact

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
	UserRepo               *repo.UserRepo
	ContactRepo            *repo.ContactRepo
	AddressRepo            *repo.AddressRepo
	ProfileRepo            *repo.ProfileRepo
	AppRepo                *repo.AppRepo
	TokenRepo              *repo.TokenRepo
	AuditLogRepo           *repo.AuditLogRepo
	JWTSigningMaterialRepo *repo.JWTSigningMaterialRepo
	IDGenerator            func(getZeroId bool) string
	SetupTestDataSource    func(t *testing.T, input RepoTestHarnessInput)
	CleanupTestDataSource  func(t *testing.T, input RepoTestHarnessInput)
}

// NOTE: The way I created the repo test harness the tests need to run
// in the order specified so that the data is there for the tests as they run.
// In heindsight I do not like this. I want to come back and improve this
// at some point, but for now it works...

func RunReposTestHarness(t *testing.T, input RepoTestHarnessInput) {
	if input.SetupTestDataSource != nil {
		input.SetupTestDataSource(t, input)
	}
	if input.CleanupTestDataSource != nil {
		defer input.CleanupTestDataSource(t, input)
	}
	setupTestHarnessData(t, input)

	// functionality tests
	t.Run("userRepo", func(t *testing.T) {
		if input.UserRepo != nil {
			testUserRepo(t, input)
		} else {
			t.Skip("no implementation for provided for userRepo")
		}
	})

	t.Run("contactRepo", func(t *testing.T) {
		if input.ContactRepo != nil {
			testContactRepo(t, input)
		} else {
			t.Skip("no implementation for provided for contactRepo")
		}
	})

	t.Run("addressRepo", func(t *testing.T) {
		if input.AddressRepo != nil {
			testAddressRepo(t, input)
		} else {
			t.Skip("no implementation for provided for addressRepo")
		}
	})

	t.Run("profileRepo", func(t *testing.T) {
		if input.ProfileRepo != nil {
			testProfileRepo(t, input)
		} else {
			t.Skip("no implementation for provided for profileRepo")
		}
	})

	t.Run("appRepo", func(t *testing.T) {
		if input.AppRepo != nil {
			testAppRepo(t, input)
		} else {
			t.Skip("no implementation for provided for appRepo")
		}
	})

	t.Run("tokenRepo", func(t *testing.T) {
		if input.TokenRepo != nil {
			testTokenRepo(t, input)
		} else {
			t.Skip("no implementation for provided for tokenRepo")
		}
	})

	t.Run("auditLogRepo", func(t *testing.T) {
		if input.AuditLogRepo != nil {
			testAuditLogRepo(t, input)
		} else {
			t.Skip("no implementation for provided for auditLogRepo")
		}
	})

	t.Run("jwtSigningMaterialRepo", func(t *testing.T) {
		if input.JWTSigningMaterialRepo != nil {
			testJWTSigningMaterialRepo(t, input)
		} else {
			t.Skip("no implementation for provided for jwtSigningMaterialRepo")
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
		t.Log(err.Error())
		t.Errorf("setup failed to add user to database: %s", err.GetErrorCode())
	}

	// create test contact
	initialTestConfirmedPrimaryContact = models.NewContact(initialTestUser.ID, "", "InitialTestUser@email.com", core.CONTACT_TYPE_EMAIL, true)
	initialTestConfirmedPrimaryContact.ConfirmedDate.Set(time.Now().UTC().Add(time.Second * -1))
	// add a test contact for the test user.
	err = (*input.ContactRepo).AddContact(context.TODO(), &initialTestConfirmedPrimaryContact, createdByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("setup failed to add contact to database: %s", err.GetErrorCode())
	}
	initialTestUnconfirmedContact = models.NewContact(initialTestUser.ID, "", "InitialTestUser2@email.com", core.CONTACT_TYPE_EMAIL, false)
	// add a test contact for the test user.
	err = (*input.ContactRepo).AddContact(context.TODO(), &initialTestUnconfirmedContact, createdByID)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("setup failed to add contact to database: %s", err.GetErrorCode())
	}

	if input.AppRepo != nil {
		// create test apps
		initialTestApp, _, err = models.NewApp(initialTestUser.ID, "test app 1", "https://my.app/callback", "https://my.app/assets/logo.png")
		if err != nil {
			t.Log(err.Error())
			t.Errorf("setup failed to create app: %s", err.GetErrorCode())
		}
		initialTestApp2, _, err = models.NewApp(initialTestUser.ID, "test app 1", "https://my.app/callback", "https://my.app/assets/logo.png")
		if err != nil {
			t.Log(err.Error())
			t.Errorf("setup failed to create app: %s", err.GetErrorCode())
		}
		// add test apps
		err = (*input.AppRepo).AddApp(context.TODO(), &initialTestApp, createdByID)
		if err != nil {
			t.Log(err.Error())
			t.Errorf("setup failed to add app to database: %s", err.GetErrorCode())
		}
		err = (*input.AppRepo).AddApp(context.TODO(), &initialTestApp2, createdByID)
		if err != nil {
			t.Log(err.Error())
			t.Errorf("setup failed to add app to database: %s", err.GetErrorCode())
		}
		// create test scopes
		initialTestAppScopes = make([]models.Scope, 0, numScopes)
		// initialTestApp2Scopes = make([]models.Scope, 0, numScopes)
		// add scopes to test app
		for i := 1; i <= numScopes; i++ {
			scope := models.NewScope(initialTestApp.ID, fmt.Sprintf("app_scope_%d", i), "scope display name", fmt.Sprintf("permissions associated with app_scope_%d", i))
			err = (*input.AppRepo).AddScope(context.TODO(), &scope, createdByID)
			if err != nil {
				t.Log(err.Error())
				t.Errorf("setup failed to add scope %d to database: %s", i, err.GetErrorCode())
			}
			initialTestAppScopes = append(initialTestAppScopes, scope)

			// scopes will be added to initialTestApp2 in the add scope test on the app repo tests.

			// scope2 := models.NewScope(initialTestApp2.ID, fmt.Sprintf("other_app_scope_%d", i), fmt.Sprintf("permissions associated with other_app_scope_%d", i))
			// err = (*input.AppRepo).AddScope(context.TODO(), &scope2, createdByID)
			// if err != nil {
			// t.Log(err.Error())
			// 	t.Errorf("setup failed to add other app scope %d to database: %s", i, err.GetErrorCode())
			// }
			// initialTestApp2Scopes = append(initialTestApp2Scopes, scope2)
		}
	}
}
