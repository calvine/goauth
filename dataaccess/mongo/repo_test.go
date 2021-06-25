package mongo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/calvine/goauth/core/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testUserRepo *userRepo
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

		t.Run("userRepo", func(t *testing.T) {
			testMongoUserRepo(t, testUserRepo)
		})

		// cleanup
		client.Database(testUserRepo.dbName).Collection(testUserRepo.collectionName).Drop(context.TODO())
	} else {
		t.Skip(SKIP_MONGO_TESTS_MESSAGE)
	}
}
