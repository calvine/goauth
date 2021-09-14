package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repo "github.com/calvine/goauth/core/repositories"
)

var (
	initialTestUser = models.User{
		PasswordHash:  "passwordhash2",
		LastLoginDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	initialTestContact = models.Contact{
		IsPrimary: true,
		Principal: "InitialTestUser@email.com",
		Type:      core.CONTACT_TYPE_EMAIL,
	}
)

type RepoTestHarnessInput struct {
	UserRepo              *repo.UserRepo
	ContactRepo           *repo.ContactRepo
	AddressRepo           *repo.AddressRepo
	ProfileRepo           *repo.ProfileRepo
	TokenRepo             *repo.TokenRepo
	AuditLogRepo          *repo.AuditLogRepo
	SetupTestDataSource   func(t *testing.T, input RepoTestHarnessInput)
	CleanupTestDataSource func(t *testing.T, input RepoTestHarnessInput)
}

func RunReposTestHarness(t *testing.T, implementationName string, input RepoTestHarnessInput) {

	setupTestHarnessData(t, input)
	if input.SetupTestDataSource != nil {
		input.SetupTestDataSource(t, input)
	}
	if input.CleanupTestDataSource != nil {
		defer input.CleanupTestDataSource(t, input)
	}

	// functionality tests
	t.Run("userRepo", func(t *testing.T) {
		if input.UserRepo != nil {
			testUserRepo(t, *input.UserRepo)
		} else {
			t.Skipf("no implementation for %s provided for userRepo", implementationName)
		}
	})

	t.Run("contactRepo", func(t *testing.T) {
		if input.ContactRepo != nil {
			testContactRepo(t, *input.ContactRepo)
		} else {
			t.Skipf("no implementation for %s provided for contactRepo", implementationName)
		}
	})

	t.Run("addressRepo", func(t *testing.T) {
		if input.AddressRepo != nil {
			// testAddressRepo(t, *input.AddressRepo)
		} else {
			t.Skipf("no implementation for %s provided for addressRepo", implementationName)
		}
	})

	t.Run("profileRepo", func(t *testing.T) {
		if input.ProfileRepo != nil {
			// testProfileRepo(t, *input.ProfileRepo)
		} else {
			t.Skipf("no implementation for %s provided for profileRepo", implementationName)
		}
	})

	t.Run("tokenRepo", func(t *testing.T) {
		if input.TokenRepo != nil {
			// testTokenRepo(t, *input.TokenRepo)
		} else {
			t.Skipf("no implementation for %s provided for tokenRepo", implementationName)
		}
	})

	t.Run("auditLogRepo", func(t *testing.T) {
		if input.AuditLogRepo != nil {
			// testAuditLogRepo(t, *input.AuditLogRepo)
		} else {
			t.Skipf("no implementation for %s provided for auditLogRepo", implementationName)
		}
	})
}

func setupTestHarnessData(t *testing.T, input RepoTestHarnessInput) {
	createdById := "test setup"
	// add a test user
	err := (*input.UserRepo).AddUser(context.TODO(), &initialTestUser, createdById)
	if err != nil {
		t.Error("setup failed to add user to database", err)
	}
	initialTestContact.UserID = initialTestUser.ID
	// add a test contact for the test user.
	err = (*input.ContactRepo).AddContact(context.TODO(), &initialTestContact, createdById)
	if err != nil {
		t.Error("setup failed to add contact to database", err)
	}
}
