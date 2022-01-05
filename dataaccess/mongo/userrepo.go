package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repomodels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ProjUserOnly = bson.M{
		"_id":                            1,
		"passwordHash":                   1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"lastLoginDate":                  1,
	}
	ProjUserWithSpecificContact = bson.M{
		"_id":                            1,
		"passwordHash":                   1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"lastLoginDate":                  1,
		"contacts.$":                     1,
	}
)

// userRepo is the repository struct for the user side of mongo db access. since other models related to users are embedded it makes sense (at least right now) to use a single struct for the related repository interfaces.
type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) userRepo {
	return userRepo{client, DB_NAME, USER_COLLECTION}
}

func NewUserRepoWithNames(client *mongo.Client, dbName, collectionName string) userRepo {
	return userRepo{client, dbName, collectionName}
}

func (userRepo) GetName() string {
	return "userRepo"
}

func (userRepo) GetType() string {
	return dataSourceType
}

func (ur userRepo) GetUserByID(ctx context.Context, id string) (models.User, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserByID", ur.GetType())
	defer span.End()
	var repoUser repomodels.RepoUser
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		rErr := coreerrors.NewFailedToParseObjectIDError(id, err, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), id)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return repoUser.ToCoreUser(), rErr
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&repoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id": id,
			}
			rErr := coreerrors.NewNoUserFoundError(fields, true)
			evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), id)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return models.User{}, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return models.User{}, rErr
	}
	user := repoUser.ToCoreUser()
	span.AddEvent("user retreived")
	return user, nil
}

func (ur userRepo) GetUserAndContactByConfirmedContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserAndContactByConfirmedContact", ur.GetType())
	defer span.End()
	var receiver struct {
		User    repomodels.RepoUser      `bson:",inline"`
		Contact []repomodels.RepoContact `bson:"contacts"`
	}
	var user models.User
	var contact models.Contact

	options := options.FindOneOptions{
		Projection: ProjUserWithSpecificContact,
	}
	filter := bson.M{
		"contacts": bson.D{
			{
				Key: "$elemMatch", Value: bson.D{
					{Key: "type", Value: contactType},
					{Key: "principal", Value: contactPrincipal},
					{Key: "confirmedDate", Value: bson.M{
						"$ne": nil,
					}},
				},
			},
		},
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"contacts.type":      contactType,
				"contacts.principal": contactPrincipal,
			}
			rErr := coreerrors.NewNoUserFoundError(fields, true)
			evtString := fmt.Sprintf("no user found with contact %s of type %s", contactPrincipal, contactType)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return user, contact, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return user, contact, rErr
	}
	user = receiver.User.ToCoreUser()
	contact = receiver.Contact[0].ToCoreContact()
	contact.UserID = user.ID
	span.AddEvent("user and contact retreived")
	return user, contact, nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetUserByPrimaryContact", ur.GetType())
	defer span.End()
	var repoUser repomodels.RepoUser
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
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"contacts.isPrimary": true,
				"contacts.type":      contactType,
				"contacts.principal": contactPrincipal,
			}
			rErr := coreerrors.NewNoUserFoundError(fields, true)
			evtString := fmt.Sprintf("no user found with primary contact %s of type %s", contactPrincipal, contactType)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return user, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return user, rErr
	}
	span.AddEvent("user retreived")
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "AddUser", ur.GetType())
	defer span.End()
	user.AuditData.CreatedByID = createdByID
	user.AuditData.CreatedOnDate = time.Now().UTC()
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).InsertOne(ctx, user, nil)
	if err != nil {
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	oid := result.InsertedID.(primitive.ObjectID)
	// oid, ok := result.InsertedID.(primitive.ObjectID)
	// if !ok {
	// 	return mongoerrors.NewMongoFailedToParseObjectID(result.InsertedID, true)
	// }
	user.ID = oid.Hex()
	span.AddEvent("user added")
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "UpdateUser", ur.GetType())
	defer span.End()
	user.AuditData.ModifiedByID.Set(modifiedByID)
	user.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	repoUser, err := repomodels.CoreUser(*user).ToRepoUser()
	if err != nil {
		evtString := fmt.Sprintf("failed to convert user to repo user: %s", err.GetErrorMessage())
		apptelemetry.SetSpanOriginalError(&span, err, evtString)
		return err
	}
	filter := bson.M{
		"_id": bson.M{
			"$eq": repoUser.ObjectID,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"passwordHash":                   repoUser.PasswordHash,
			"consecutiveFailedLoginAttempts": repoUser.ConsecutiveFailedLoginAttempts,
			"lockedOutUntil":                 repoUser.LockedOutUntil.GetPointerCopy(),
			"lastLoginDate":                  repoUser.LastLoginDate.GetPointerCopy(),
			"modifiedById":                   repoUser.AuditData.ModifiedByID.GetPointerCopy(),
			"modifiedOnDate":                 repoUser.AuditData.ModifiedOnDate.GetPointerCopy(),
		},
	}
	result, updateErr := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update)
	if updateErr != nil {
		rErr := coreerrors.NewRepoQueryFailedError(updateErr, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id": user.ID,
		}
		rErr := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with id: %s", user.ID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	span.AddEvent("user updated")
	return nil
}
