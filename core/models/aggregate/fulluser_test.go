package aggregate

import (
	"reflect"
	"testing"

	"github.com/calvine/goauth/core/models"
)

func TestNewFullUserWithData(t *testing.T) {
	emptyFullUser := FullUser{}
	user := models.User{}
	profile := models.Profile{}
	fullUser := NewFullUserWithData(user, nil, nil, &profile)
	if reflect.DeepEqual(fullUser, emptyFullUser) {
		t.Error("expected full user to be populated", fullUser)
	}
	if &profile == &fullUser.Profile {
		t.Error("memory addresses for profile and fullUser.Profile should not be the same", &profile, &fullUser.Profile)
	}
}
