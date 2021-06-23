package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/calvine/goauth/models/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("unable to find user with given id")
)

type UserRepo interface {
	GetUserById(ctx context.Context, id string) (core.User, error)
	AddUser(ctx context.Context, user *core.User, createdById string) error
	UpdateUser(ctx context.Context, user *core.User, modifiedById string) error
}

type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) *userRepo {
	return &userRepo{client, DB_NAME, USER_COLLECTION}
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

func (ur userRepo) AddUser(ctx context.Context, user *core.User, createdById string) error {
	_, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).InsertOne(ctx, user, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *core.User, modifiedById string) error {
	user.ModifiedByID.Set(modifiedById)
	user.ModifiedOnDate.Set(time.Now().UTC())
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, bson.M{}, user)
	if result.ModifiedCount == 0 {
		return ErrUserNotFound
	}
	if err != nil {
		return nil
	}
	return nil
}
