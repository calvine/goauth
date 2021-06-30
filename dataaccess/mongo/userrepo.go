package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/calvine/goauth/core/models"
	repoModels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	emptyUser = models.User{}

	ErrUserNotFound          = errors.New("unable to find user with given id")
	ErrFailedToParseObjectId = errors.New("failed to parse object id")

	ProjUserOnly = bson.M{
		"_id":                            1,
		"passwordHash":                   1,
		"salt":                           1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"lastLoginDate":                  1,
	}
)

// userRepo is the repository struct for the user side of mongo db access. since other models related to users are embedded it makes sense (at least right now) to use a single struct for the related repository interfaces.
type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) *userRepo {
	return &userRepo{client, DB_NAME, USER_COLLECTION}
}

func NewUserRepoWithNames(client *mongo.Client, dbName, collectionName string) *userRepo {
	return &userRepo{client, dbName, collectionName}
}

func (ur userRepo) GetUserById(ctx context.Context, id string) (models.User, error) {
	var repoUser repoModels.RepoUser
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repoUser.ToCoreUser(), err
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&repoUser)
	user := repoUser.ToCoreUser()
	if err != nil {
		return user, err
	}
	if user == emptyUser {
		return user, ErrUserNotFound
	}
	return user, nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, error) {
	var repoUser repoModels.RepoUser
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	filter := bson.M{
		"contacts": bson.D{
			{
				Key: "$elemMatch", Value: bson.D{
					{Key: "isPrimary", Value: true},
					{Key: "type", Value: contactType},
					{Key: "principal", Value: contactPrincipal},
				},
			},
		},
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&repoUser)
	user := repoUser.ToCoreUser()
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdById string) error {
	user.AuditData.CreatedByID = createdById
	user.AuditData.CreatedOnDate = time.Now().UTC()
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).InsertOne(ctx, user, nil)
	if err != nil {
		return err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return ErrFailedToParseObjectId
	}
	user.ID = oid.Hex()
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedById string) error {
	user.AuditData.ModifiedByID.Set(modifiedById)
	user.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	repoUser, err := repoModels.CoreUser(*user).ToRepoUser()
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": bson.M{
			"$eq": repoUser.ObjectId,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"passwordHash":                   repoUser.PasswordHash,
			"salt":                           repoUser.Salt,
			"consecutiveFailedLoginAttempts": repoUser.ConsecutiveFailedLoginAttempts,
			"lockedOutUntil":                 repoUser.LockedOutUntil.GetPointerCopy(),
			"lastLoginDate":                  repoUser.LastLoginDate.GetPointerCopy(),
			"modifiedById":                   repoUser.AuditData.ModifiedByID.GetPointerCopy(),
			"modifiedOnDate":                 repoUser.AuditData.ModifiedOnDate.GetPointerCopy(),
		},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
