package mongo

import (
	"context"
	"time"

	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Contacts []models.Contact `bson:"contacts"`
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
		return emptyContact, err
	}
	if len(receiver.Contacts) == 0 {
		return emptyContact, ErrUserNotFound
	}
	// TODO: need to make sure business logic exists to ensure that there is only 1 primary contact...
	contact := receiver.Contacts[0]
	contact.UserId = userId
	return contact, nil
}

// TODO: finish implementing

func (ur *userRepo) GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, error) {
	var receiver struct {
		Contacts []models.Contact `bson:"contacts"`
	}
	options := options.FindOneOptions{
		Projection: bson.D{
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
		return nil, err
	}
	return receiver.Contacts, nil
}

func (ur *userRepo) GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, error) {
	var receiver struct {
		id      primitive.ObjectID `bson:"_id"`
		contact models.Contact     `bson:"contacts"`
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
		return emptyContact, err
	}
	return receiver.contact, nil
}

func (ur *userRepo) AddContact(ctx context.Context, contact *models.Contact, createdById string) error {
	contact.AuditData.CreatedById = createdById
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	contact.Id = "new id / uuid?"
	oid, err := primitive.ObjectIDFromHex(contact.UserId)
	if err != nil {
		return ErrFailedToParseObjectId
	}
	update := bson.M{
		"$push": bson.M{
			"contacts": contact,
		},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateByID(ctx, oid, update) //(ctx, contact, nil)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return ErrUserNotFound
	}
	return nil
}

// func (ur *userRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) error {
// TODO: Use array filters?
//  contact.AuditData.ModifiedById = modifiedById
// 	contact.AuditData.ModifiedOnDate = time.Now().UTC()
// 	contact.Id = "new id / uuid?"
// 	oid, err := primitive.ObjectIDFromHex(contact.UserId)
// 	if err != nil {
// 		// TODO: specific error here?
// 		return err
// 	}
// 	arrayFilters := options.ArrayFilters{
// 		Filters: bson.A{},
// 	}
// 	options := options.UpdateOptions{
// 		ArrayFilters: &arrayFilters,
// 		Upsert:       false,
// 	}
// 	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateByID(ctx, oid, nil, &options) //(ctx, contact, nil)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
