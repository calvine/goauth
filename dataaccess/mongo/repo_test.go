package mongo

import "fmt"

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
