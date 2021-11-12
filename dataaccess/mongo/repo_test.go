package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/calvine/goauth/core/normalization"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/dataaccess/internal/repotest"

	"github.com/calvine/goauth/core/utilities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	value, exists := os.LookupEnv(ENV_RUN_MONGO_TESTS)
	shouldRun, _ := normalization.ReadBoolValue(value, true)
	if exists && shouldRun {
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
		testUserRepo := NewUserRepoWithNames(client, "test_goauth", USER_COLLECTION)
		var userRepo repo.UserRepo = testUserRepo
		var contactRepo repo.ContactRepo = testUserRepo
		cleanUpDataSource := func(t *testing.T, _ repotest.RepoTestHarnessInput) {
			err := testUserRepo.mongoClient.Database(testUserRepo.dbName).Collection(testUserRepo.collectionName).Drop(context.TODO())
			if err != nil {
				t.Error("failed to cleanup database", err)
			}
		}
		testHarnessInput := repotest.RepoTestHarnessInput{
			UserRepo:            &userRepo,
			ContactRepo:         &contactRepo,
			SetupTestDataSource: cleanUpDataSource,
			IDGenerator: func(getZeroId bool) string {
				if getZeroId {
					return primitive.NilObjectID.Hex()
				}
				return primitive.NewObjectID().Hex()
			},
		}
		repotest.RunReposTestHarness(t, "mongodb", testHarnessInput)
	} else {
		t.Skip(SKIP_MONGO_TESTS_MESSAGE)
	}
}
