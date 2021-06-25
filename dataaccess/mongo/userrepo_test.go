package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testUser1 = models.User{
		Id:           "uniqueId1",
		PasswordHash: "passwordhash1",
		Salt:         "salt1",
	}
	testUser2 = models.User{
		Id:           "uniqueId2",
		PasswordHash: "passwordhash2",
		Salt:         "salt2",
	}
)

func TestMongoUserRepo(t *testing.T) {
	// setup code for mongo user repo tests.
	_, exists := os.LookupEnv(ENV_RUN_MONGO_TESTS)
	if !exists {
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

		// functionality tests
		t.Run("userRepo.AddUser", func(t *testing.T) {
			testAddUser(t, testUserRepo)
		})
		t.Run("userRepo.UpdateUser", func(t *testing.T) {
			testUpdateUser(t, testUserRepo)
		})
		t.Run("userRepo.GetUserById", func(t *testing.T) {
			t.Fail()
		})
		t.Run("userRepo.GetUserByPrimaryContact", func(t *testing.T) {
			t.Fail()
		})

		// cleanup
		client.Database(testUserRepo.dbName).Collection(testUserRepo.collectionName).Drop(context.TODO())
	} else {
		t.Skip(SKIP_MONGO_TESTS_MESSAGE)
	}
}

func testAddUser(t *testing.T, userRepo *userRepo) {
	createdById := "test1"

	err := userRepo.AddUser(context.TODO(), &testUser1, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser1.AuditData.CreatedById != createdById {
		t.Error("failed to set the users CreatedByID to the right value", testUser1.AuditData.CreatedById, createdById)
	}

	createdById = "test2"

	err = userRepo.AddUser(context.TODO(), &testUser2, createdById)
	if err != nil {
		t.Error("failed to add user to database", err)
	}
	if testUser2.AuditData.CreatedById != createdById {
		t.Error("failed to set the users CreatedByID to the right value", testUser2.AuditData.CreatedById, createdById)
	}
}

func testUpdateUser(t *testing.T, userRepo *userRepo) {
	preUpdateDate := time.Now().UTC()
	newPasswordHash := "another secure password hash"
	newSalt := "change password = change salt"
	testUser1.PasswordHash = newPasswordHash
	testUser1.Salt = newSalt
	err := userRepo.UpdateUser(context.TODO(), &testUser1, testUser1.Id)
	if err != nil {
		t.Error("failed to update user", err)
	}
	if !testUser1.AuditData.ModifiedOnDate.Value.After(preUpdateDate) {
		t.Error("ModifiedOnDate should be after the preUpdate for test", preUpdateDate, testUser1.AuditData.ModifiedOnDate)
	}
	retreivedUser, err := userRepo.GetUserById(context.TODO(), testUser1.Id)
	if err != nil {
		t.Error("failed to retreive updated user to check that fields were updated.", err)
	}
	if retreivedUser.PasswordHash != newPasswordHash {
		t.Error("password hash did not update.", retreivedUser.PasswordHash, newPasswordHash)
	}
	if retreivedUser.Salt != newSalt {
		t.Error("salt did not update.", retreivedUser.Salt, newSalt)
	}
}
