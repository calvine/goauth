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
)

var (
	initialTestUser models.User

	initialTestContact models.Contact

	initialTestApp models.App

	initialTestAppScopes []models.Scope
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
	createdByID := "test setup"
	// create test user
	initialTestUser = models.NewUser()
	initialTestUser.PasswordHash = "passwordhash2"
	initialTestUser.LastLoginDate = nullable.NullableTime{HasValue: true, Value: time.Now().UTC()}
	// add a test user
	err := (*input.UserRepo).AddUser(context.TODO(), &initialTestUser, createdByID)
	if err != nil {
		t.Error("setup failed to add user to database", err)
	}

	// create test contact
	initialTestContact = models.NewContact(initialTestUser.ID, "", "InitialTestUser@email.com", core.CONTACT_TYPE_EMAIL, true)
	// add a test contact for the test user.
	err = (*input.ContactRepo).AddContact(context.TODO(), &initialTestContact, createdByID)
	if err != nil {
		t.Error("setup failed to add contact to database", err)
	}

	// create test app

	// add a test app

	// create test scopes

	// add scopes to test app
}
