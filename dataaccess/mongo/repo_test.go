package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/goauth/core/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testUserRepo *userRepo

	initialTestUser = models.User{
		PasswordHash:  "passwordhash2",
		Salt:          "salt2",
		LastLoginDate: nullable.NullableTime{HasValue: true, Value: time.Now().UTC()},
	}

	initialTestContact = models.Contact{
		ConfirmationCode: nullable.NullableString{HasValue: true, Value: "abc123"},
		IsPrimary:        true,
		Principal:        "InitialTestUser@email.com",
		Type:             core.CONTACT_TYPE_EMAIL,
	}
)

const (
	ENV_RUN_MONGO_TESTS              = "GOAUTH_RUN_MONGO_TESTS"
	ENV_MONGO_TEST_CONNECTION_STRING = "GOAUTH_MONGO_TEST_CONNECTION_STRING"

	DEFAULT_TEST_MONGO_CONNECTION_STRING = "mongodb://root:password@localhost:27017/?authSource=admin&readPreference=primary&ssl=false"
)

var (
	SKIP_MONGO_TESTS_MESSAGE = fmt.Sprintf("skipping mongo tests because env var %s was not set", ENV_RUN_MONGO_TESTS)
)

func TestMongoRepos(t *testing.T) {
	_, exists := os.LookupEnv(ENV_RUN_MONGO_TESTS)
	// TODO: remove ! so this chek works properly
	if !exists {
		// setup code for mongo user repo tests.
		connectionString := utilities.GetEnv(ENV_MONGO_TEST_CONNECTION_STRING, DEFAULT_TEST_MONGO_CONNECTION_STRING)
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
		defer client.Disconnect(context.TODO())
		if err != nil {
			t.Error("failed to connect to mongo server", err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			t.Error("failed to ping mongo server before test", err)
		}
		testUserRepo = NewUserRepoWithNames(client, "test_goauth", USER_COLLECTION)
		// cleanup previous test data if it exists
		cleanupTestDatabase(testUserRepo)
		// set up initial data for testing
		setupTestData(t, testUserRepo)

		// functionality tests
		t.Run("userRepo", func(t *testing.T) {
			testMongoUserRepo(t, testUserRepo)
		})

		t.Run("contactRepo", func(t *testing.T) {
			testMongoContactRepo(t, testUserRepo)
		})
	} else {
		t.Skip(SKIP_MONGO_TESTS_MESSAGE)
	}
}

func cleanupTestDatabase(userRepo *userRepo) error {
	return userRepo.mongoClient.Database(testUserRepo.dbName).Collection(testUserRepo.collectionName).Drop(context.TODO())
}

func setupTestData(t *testing.T, userRepo *userRepo) {
	createdById := "test setup"
	// add a test user
	err := userRepo.AddUser(context.TODO(), &initialTestUser, createdById)
	if err != nil {
		t.Error("setup failed to add user to database", err)
	}
	initialTestContact.UserId = initialTestUser.Id
	// add a test contact for the test user.
	err = userRepo.AddContact(context.TODO(), &initialTestContact, createdById)
	if err != nil {
		t.Error("setup failed to add contact to database", err)
	}
}
