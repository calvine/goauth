package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound = errors.New("unable to find user with given id")

	ProjUserOnly = bson.M{
		"id":                             1,
		"password":                       1,
		"salt":                           1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"LastLoginDate":                  1,
	}
)

type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) *userRepo {
	return &userRepo{client, DB_NAME, USER_COLLECTION}
}

func (ur userRepo) GetUserById(ctx context.Context, id string) (models.User, error) {
	var user models.User
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, bson.M{"userId": id}, &options).Decode(&user)
	if err != nil {
		return user, err
	}
	if (user == models.User{}) {
		return user, ErrUserNotFound
	}
	return user, nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactPrincipalType, contactPrincipal string) (models.User, error) {
	var user models.User
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	result := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, bson.M{}, &options)
	err := result.Err()
	if err != nil {
		return user, err
	}
	err = result.Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdById string) error {
	_, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).InsertOne(ctx, user, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedById string) error {
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
