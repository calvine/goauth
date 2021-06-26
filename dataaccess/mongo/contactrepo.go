package mongo

import (
	"context"

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
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, bson.M{"_id": oid}, &options).Decode(&contact)
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

// func (ur *userRepo) GetContactsByUserId(ctx context.Context, userId string) ([]models.Contact, error) {
// 	var contacts []models.Contact
// 	options := options.FindOneOptions{
// 		Projection: ProjContactOnly,
// 	}
// 	oid, err := primitive.ObjectIDFromHex(userId)
// 	if err != nil {
// 		return contacts, err
// 	}
// 	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, bson.M{"_id": oid}, &options).Decode(&contacts)

// 	if err != nil {
// 		return contacts, err
// 	}
// 	if contacts == emptyContact {
// 		return contacts, ErrUserNotFound
// 	}
// 	return contacts, nil
// }

// func (ur *userRepo) GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, error) {

// }

// func (ur *userRepo) AddContact(ctx context.Context, contact *models.Contact, createdById string) error {

// }

// func (ur *userRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedById string) error {

// }
