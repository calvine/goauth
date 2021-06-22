package mongo

import (
	"context"
	"errors"

	"github.com/calvine/goauth/models/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("unable to find user with given id")
)

type UserRepo interface {
	GetUserById(ctx context.Context, id string) (core.User, error)
	AddUser(ctx context.Context, user *core.User) error
	UpdateUser(ctx context.Context, user *core.User) error
}

type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) *userRepo {
	return &userRepo{client, "", ""}
}

func (ur userRepo) GetUserById(ctx context.Context, id string) (core.User, error) {
	var user core.User
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, bson.M{"userId": id}).Decode(&user)
	if err != nil {
		return user, err
	}
	if (user == core.User{}) {
		return user, ErrUserNotFound
	}
	return user, nil
}
