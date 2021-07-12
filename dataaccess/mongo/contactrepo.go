package mongo

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repoModels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	emptyContact = models.Contact{}

	ProjContactOnly = bson.M{
		"id":               1,
		"name":             1,
		"principal":        1,
		"type":             1,
		"isPrimary":        1,
		"confirmationCode": 1,
		"confirmedDate":    1,
	}
)

func (ur *userRepo) GetPrimaryContactByUserId(ctx context.Context, userId string) (models.Contact, error) {
	var receiver struct {
		Contacts []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{}
	options.SetProjection(bson.D{
		{Key: "_id", Value: 0},
		{Key: "contacts.$", Value: 1},
	})
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return emptyContact, ErrFailedToParseObjectId
	}
	filter := bson.M{
		"$and": bson.A{
			bson.M{"_id": oid},
			bson.M{"contacts.isPrimary": true},
		},
	}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id":                userId,
				"contacts.isPrimary": true,
			}
			return emptyContact, coreerrors.NewRepoNoContactFoundErrorWithFields(fields, true)
		}
		return emptyContact, coreerrors.NewRepoQueryFailed(err, true)
	}
	// TODO: need to make sure business logic exists to ensure that there is only 1 primary contact...
	contact := receiver.Contacts[0].ToCoreContact()
	contact.UserID = userId
	return contact, nil
}

// TODO: finish implementing

func (ur *userRepo) GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, error) {
	var receiver struct {
		Contacts []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{
		Projection: bson.D{
			{Key: "_id", Value: 0},
			{Key: "contacts", Value: 1},
		},
	}
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, ErrFailedToParseObjectId
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, coreerrors.NewRepoNoContactFoundError("_id", userId, true)
		}
		return nil, coreerrors.NewRepoQueryFailed(err, true)
	}
	contacts := make([]models.Contact, len(receiver.Contacts))
	for index, contact := range receiver.Contacts {
		contact.UserID = userId
		contacts[index] = contact.ToCoreContact()
	}
	return contacts, nil
}

// TODO: Figure out why decoding confirmation code is failing...
func (ur *userRepo) GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, error) {
	var receiver struct {
		UserId  primitive.ObjectID       `bson:"_id"`
		Contact []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{
		Projection: bson.D{
			{Key: " _id", Value: 1},
			{Key: "contacts.$", Value: 1},
		},
	}
	filter := bson.M{
		"contacts.confirmationCode": confirmationCode,
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return emptyContact, coreerrors.NewRepoNoContactFoundError("contacts.confirmationCode", confirmationCode, true)
		}
		return emptyContact, coreerrors.NewRepoQueryFailed(err, true)
	}
	receiver.Contact[0].UserID = receiver.UserId.Hex()
	return receiver.Contact[0].ToCoreContact(), nil
}

func (ur *userRepo) AddContact(ctx context.Context, contact *models.Contact, createdById string) error {
	contact.AuditData.CreatedByID = createdById
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	contact.ID = primitive.NewObjectID().Hex()
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		return ErrFailedToParseObjectId
	}
	repoContact, err := repoModels.CoreContact(*contact).ToRepoContact()
	if err != nil {
		return err
	}
	update := bson.M{
		"$push": bson.M{
			"contacts": bson.D{
				{Key: "id", Value: repoContact.ObjectId},
				{Key: "name", Value: repoContact.CoreContact.Name.GetPointerCopy()},
				{Key: "principal", Value: repoContact.CoreContact.Principal},
				{Key: "type", Value: repoContact.CoreContact.Type},
				{Key: "isPrimary", Value: repoContact.CoreContact.IsPrimary},
				{Key: "confirmationCode", Value: repoContact.CoreContact.ConfirmationCode.GetPointerCopy()},
				{Key: "confirmedDate", Value: repoContact.CoreContact.ConfirmedDate.GetPointerCopy()},
				{Key: "createdById", Value: repoContact.CoreContact.AuditData.CreatedByID},
				{Key: "createdOnDate", Value: repoContact.CoreContact.AuditData.CreatedOnDate},
				{Key: "modifiedById", Value: nil},
				{Key: "modifiedOnDate", Value: nil},
			},
		},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateByID(ctx, oid, update) //(ctx, contact, nil)
	if err != nil {
		return coreerrors.NewRepoQueryFailed(err, true)
	}
	if result.ModifiedCount == 0 {
		// TODO: replace with rich error.
		return ErrUserNotFound
	}
	return nil
}

func (ur *userRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) error {
	contact.AuditData.ModifiedByID = nullable.NullableString{}
	contact.AuditData.ModifiedByID.Set(modifiedById)
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{}
	contact.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	contactID, err := primitive.ObjectIDFromHex(contact.ID)
	if err != nil {
		// TODO: specific error here?
		return err
	}
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		// TODO: specific error here?
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "contacts.id", Value: contactID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "contacts.$.name", Value: contact.Name.GetPointerCopy()},
			{Key: "contacts.$.principal", Value: contact.Principal},
			{Key: "contacts.$.type", Value: contact.Type},
			{Key: "contacts.$.isPrimary", Value: contact.IsPrimary},
			{Key: "contacts.$.confirmationCode", Value: contact.ConfirmationCode.GetPointerCopy()},
			{Key: "contacts.$.confirmedDate", Value: contact.ConfirmedDate.GetPointerCopy()},
			{Key: "contacts.$.modifiedById", Value: contact.AuditData.ModifiedByID.GetPointerCopy()},
			{Key: "contacts.$.modifiedOnDate", Value: contact.AuditData.ModifiedOnDate.GetPointerCopy()},
		}},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update, nil) //(ctx, contact, nil)
	if err != nil {
		return coreerrors.NewRepoQueryFailed(err, true)
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id":        contact.UserID,
			"contact.id": contact.ID,
		}
		return coreerrors.NewRepoNoContactFoundErrorWithFields(fields, true)
	}
	return nil
}

func (ur *userRepo) ConfirmContact(ctx context.Context, confirmationCode, modifiedById string) error {
	now := time.Now().UTC()
	filter := bson.D{
		{Key: "contacts.confirmationCode", Value: confirmationCode},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "contacts.$.confirmationCode", Value: nil},
			{Key: "contacts.$.confirmedDate", Value: now},
			{Key: "contacts.$.modifiedById", Value: modifiedById},
			{Key: "contacts.$.modifiedOnDate", Value: now},
		}},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update, nil)
	if err != nil {
		return coreerrors.NewRepoQueryFailed(err, true)
	}
	if result.ModifiedCount == 0 {
		return coreerrors.NewRepoNoContactFoundError("contacts.confirmationCode", confirmationCode, true)
	}
	return nil
}
