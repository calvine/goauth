package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	emptyContact = models.Contact{}

	ProjContactOnly = bson.M{
		"id":            1,
		"name":          1,
		"principal":     1,
		"type":          1,
		"isPrimary":     1,
		"confirmCode":   1,
		"confirmedDate": 1,
	}
)

func (ur *userRepo) GetPrimaryContactByUserId(ctx context.Context, userId string) (models.Contact, error) {
	var contact models.Contact
	options := options.FindOneOptions{
		Projection: ProjContactOnly,
	}
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return contact, err
	}
	filter := bson.M{
		"$and": bson.A{
			bson.M{"_id": oid},
			bson.M{"contacts.isPrimary": true},
		},
	}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&contact)
	contact.UserId = userId
	if err != nil {
		return contact, err
	}
	if contact == emptyContact {
		return contact, ErrUserNotFound
	}
	return contact, nil
}

// TODO: finish implementing

func (ur *userRepo) GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, error) {
	var contacts []models.Contact
	options := options.FindOneOptions{
		Projection: ProjContactOnly,
	}
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return contacts, err
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&contacts)

	if err != nil {
		return contacts, err
	}
	// if contacts == emptyContact {
	// 	return contacts, ErrUserNotFound
	// }
	return contacts, nil
}

// func (ur *userRepo) GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, error) {

// }

func (ur *userRepo) AddContact(ctx context.Context, contact *models.Contact, createdById string) error {
	contact.AuditData.CreatedById = createdById
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	contact.Id = "new id / uuid?"
	oid, err := primitive.ObjectIDFromHex(contact.UserId)
	if err != nil {
		// TODO: specific error here?
		return err
	}
	update := bson.M{
		"$push": bson.D{
			{"contacts", contact},
		},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateByID(ctx, oid, update) //(ctx, contact, nil)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		// TODO: add specfic error here.
		return errors.New("no document updated")
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
