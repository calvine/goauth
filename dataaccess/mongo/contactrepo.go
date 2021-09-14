package mongo

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	repoModels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		"confirmedDate": 1,
	}
)

func (ur userRepo) GetContactByID(ctx context.Context, id string) (models.Contact, errors.RichError) {
	var receiver struct {
		UserID  primitive.ObjectID       `bson:"_id"`
		Contact []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{}
	options.SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "contacts.$", Value: 1},
	})
	contactOid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return emptyContact, coreerrors.NewFailedToParseObjectIDError(id, err, true)
	}
	filter := bson.D{
		{Key: "contacts.id", Value: contactOid},
	}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"contact.id": id,
			}
			return emptyContact, coreerrors.NewNoContactFoundError(fields, true)
		}
		return emptyContact, coreerrors.NewRepoQueryFailedError(err, true)
	}
	receiver.Contact[0].UserID = receiver.UserID.Hex()
	return receiver.Contact[0].ToCoreContact(), nil
}

func (ur userRepo) GetPrimaryContactByUserID(ctx context.Context, userID string) (models.Contact, errors.RichError) {
	var receiver struct {
		Contacts []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{}
	options.SetProjection(bson.D{
		{Key: "_id", Value: 0},
		{Key: "contacts.$", Value: 1},
	})
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return emptyContact, coreerrors.NewFailedToParseObjectIDError(userID, err, true)
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
				"_id":                userID,
				"contacts.isPrimary": true,
			}
			return emptyContact, coreerrors.NewNoContactFoundError(fields, true)
		}
		return emptyContact, coreerrors.NewRepoQueryFailedError(err, true)
	}
	// TODO: need to make sure business logic exists to ensure that there is only 1 primary contact...
	contact := receiver.Contacts[0].ToCoreContact()
	contact.UserID = userID
	return contact, nil
}

// TODO: finish implementing

func (ur userRepo) GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError) {
	var receiver struct {
		Contacts []repoModels.RepoContact `bson:"contacts"`
	}
	options := options.FindOneOptions{
		Projection: bson.D{
			{Key: "_id", Value: 0},
			{Key: "contacts", Value: 1},
		},
	}
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, coreerrors.NewFailedToParseObjectIDError(userID, err, true)
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id": userID,
			}
			return nil, coreerrors.NewNoContactFoundError(fields, true)
		}
		return nil, coreerrors.NewRepoQueryFailedError(err, true)
	}
	contacts := make([]models.Contact, len(receiver.Contacts))
	for index, contact := range receiver.Contacts {
		contact.UserID = userID
		contacts[index] = contact.ToCoreContact()
	}
	return contacts, nil
}

// func (ur userRepo) GetContactByConfirmationCode(ctx context.Context, confirmationCode string) (models.Contact, errors.RichError) {
// 	var receiver struct {
// 		UserID  primitive.ObjectID       `bson:"_id"`
// 		Contact []repoModels.RepoContact `bson:"contacts"`
// 	}
// 	options := options.FindOneOptions{
// 		Projection: bson.D{
// 			{Key: " _id", Value: 1},
// 			{Key: "contacts.$", Value: 1},
// 		},
// 	}
// 	filter := bson.M{
// 		"contacts.confirmationCode": confirmationCode,
// 	}
// 	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			fields := map[string]interface{}{
// 				"contacts.confirmationCode": confirmationCode,
// 			}
// 			return emptyContact, coreerrors.NewNoContactFoundError(fields, true)
// 		}
// 		return emptyContact, coreerrors.NewRepoQueryFailedError(err, true)
// 	}
// 	receiver.Contact[0].UserID = receiver.UserID.Hex()
// 	return receiver.Contact[0].ToCoreContact(), nil
// }

func (ur userRepo) AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError {
	contact.AuditData.CreatedByID = createdByID
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	contact.ID = primitive.NewObjectID().Hex()
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		return coreerrors.NewFailedToParseObjectIDError(contact.UserID, err, true)
	}
	repoContact, convertErr := repoModels.CoreContact(*contact).ToRepoContact()
	if err != nil {
		return convertErr
	}
	update := bson.M{
		"$push": bson.M{
			"contacts": bson.D{
				{Key: "id", Value: repoContact.ObjectID},
				{Key: "name", Value: repoContact.CoreContact.Name.GetPointerCopy()},
				{Key: "principal", Value: repoContact.CoreContact.Principal},
				{Key: "type", Value: repoContact.CoreContact.Type},
				{Key: "isPrimary", Value: repoContact.CoreContact.IsPrimary},
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
		return coreerrors.NewRepoQueryFailedError(err, true)
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id": contact.UserID,
		}
		return coreerrors.NewNoUserFoundError(fields, true)
	}
	return nil
}

func (ur userRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError {
	contact.AuditData.ModifiedByID = nullable.NullableString{}
	contact.AuditData.ModifiedByID.Set(modifiedByID)
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{}
	contact.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	contactID, err := primitive.ObjectIDFromHex(contact.ID)
	if err != nil {
		// TODO: specific error here?
		return coreerrors.NewFailedToParseObjectIDError(contact.ID, err, true)
	}
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		// TODO: specific error here?
		return coreerrors.NewFailedToParseObjectIDError(contact.UserID, err, true)
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
			{Key: "contacts.$.confirmedDate", Value: contact.ConfirmedDate.GetPointerCopy()},
			{Key: "contacts.$.modifiedById", Value: contact.AuditData.ModifiedByID.GetPointerCopy()},
			{Key: "contacts.$.modifiedOnDate", Value: contact.AuditData.ModifiedOnDate.GetPointerCopy()},
		}},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update, nil) //(ctx, contact, nil)
	if err != nil {
		return coreerrors.NewRepoQueryFailedError(err, true)
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id":        contact.UserID,
			"contact.id": contact.ID,
		}
		return coreerrors.NewNoContactFoundError(fields, true)
	}
	return nil
}

// func (ur userRepo) ConfirmContact(ctx context.Context, confirmationCode, modifiedByID string) errors.RichError {
// 	now := time.Now().UTC()
// 	filter := bson.D{
// 		{Key: "contacts.confirmationCode", Value: confirmationCode},
// 	}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "contacts.$.confirmationCode", Value: nil},
// 			{Key: "contacts.$.confirmedDate", Value: now},
// 			{Key: "contacts.$.modifiedById", Value: modifiedByID},
// 			{Key: "contacts.$.modifiedOnDate", Value: now},
// 		}},
// 	}
// 	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update, nil)
// 	if err != nil {
// 		return coreerrors.NewRepoQueryFailedError(err, true)
// 	}
// 	if result.ModifiedCount == 0 {
// 		fields := map[string]interface{}{
// 			"contacts.confirmationCode": confirmationCode,
// 		}
// 		return coreerrors.NewNoContactFoundError(fields, true)
// 	}
// 	return nil
// }
