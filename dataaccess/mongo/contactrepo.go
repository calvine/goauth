package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/calvine/goauth/core/apptelemetry"
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
		"rawPrincipal":  1,
		"type":          1,
		"isPrimary":     1,
		"confirmedDate": 1,
	}
)

func (ur userRepo) GetContactByID(ctx context.Context, id string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetContactByID", ur.GetType())
	defer span.End()
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
		rErr := coreerrors.NewFailedToParseObjectIDError(id, err, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), id)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return emptyContact, rErr
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
			rErr := coreerrors.NewNoContactFoundError(fields, true)
			evtString := fmt.Sprintf("no contact found with id: %s", id)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return emptyContact, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return emptyContact, rErr
	}
	receiver.Contact[0].UserID = receiver.UserID.Hex()
	span.AddEvent("contact retreived")
	return receiver.Contact[0].ToCoreContact(), nil
}

func (ur userRepo) GetPrimaryContactByUserID(ctx context.Context, userID string, contactType string) (models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetPrimaryContactByUserID", ur.GetType())
	defer span.End()
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
		rErr := coreerrors.NewFailedToParseObjectIDError(userID, err, true)
		evtString := fmt.Sprintf("%s: %s", rErr.GetErrorMessage(), userID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return emptyContact, rErr
	}
	filter := bson.M{
		"_id": oid,
		"contacts": bson.D{
			{
				Key: "$elemMatch", Value: bson.D{
					{Key: "type", Value: contactType},
					{Key: "isPrimary", Value: true},
				},
			},
		},
	}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id":                userID,
				"contacts.isPrimary": true,
				"contacts.type":      contactType,
			}
			rErr := coreerrors.NewNoContactFoundError(fields, true)
			evtString := fmt.Sprintf("no primary contact found for user id: %s", userID)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return emptyContact, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return emptyContact, rErr
	}
	// TODO: need to make sure business logic exists to ensure that there is only 1 primary contact...
	contact := receiver.Contacts[0].ToCoreContact()
	contact.UserID = userID
	span.AddEvent("primary contact retreived")
	return contact, nil
}

func (ur userRepo) GetContactsByUserID(ctx context.Context, userID string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetContactsByUserID", ur.GetType())
	defer span.End()
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
		rErr := coreerrors.NewFailedToParseObjectIDError(userID, err, true)
		evtString := fmt.Sprintf("%s user id: %s", rErr.GetErrorMessage(), userID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id": userID,
			}
			rErr := coreerrors.NewNoUserFoundError(fields, true)
			evtString := fmt.Sprintf("no contact found for user id: %s", userID)
			apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			return nil, rErr
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	contacts := make([]models.Contact, len(receiver.Contacts))
	for index, contact := range receiver.Contacts {
		contact.UserID = userID
		contacts[index] = contact.ToCoreContact()
	}
	span.AddEvent("contacts retreived")
	return contacts, nil
}

func (ur userRepo) GetContactsByUserIDAndType(ctx context.Context, userID string, contactType string) ([]models.Contact, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetContactsByUserID", ur.GetType())
	defer span.End()
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
		rErr := coreerrors.NewFailedToParseObjectIDError(userID, err, true)
		evtString := fmt.Sprintf("%s user id: %s", rErr.GetErrorMessage(), userID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	filter := bson.M{
		"_id":           oid,
		"contacts.type": contactType,
	}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// fields := map[string]interface{}{
			// 	"_id": userID,
			// }
			// rErr := coreerrors.NewNoUserFoundError(fields, true)
			// evtString := fmt.Sprintf("no contact found for user id: %s", userID)
			// apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
			span.AddEvent("no contacts of type found")
			return []models.Contact{}, nil
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return nil, rErr
	}
	contacts := make([]models.Contact, 0, len(receiver.Contacts))
	for _, contact := range receiver.Contacts {
		if contact.Type == contactType {
			contact.UserID = userID
			contacts = append(contacts, contact.ToCoreContact())
		}
	}
	span.AddEvent("contacts of type retreived")
	return contacts, nil
}

func (ur userRepo) AddContact(ctx context.Context, contact *models.Contact, createdByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "AddContact", ur.GetType())
	defer span.End()
	contact.AuditData.CreatedByID = createdByID
	contact.AuditData.CreatedOnDate = time.Now().UTC()
	contact.ID = primitive.NewObjectID().Hex()
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		rErr := coreerrors.NewFailedToParseObjectIDError(contact.UserID, err, true)
		evtString := fmt.Sprintf("%s user id: %s", rErr.GetErrorMessage(), contact.UserID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
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
				{Key: "rawPrincipal", Value: repoContact.CoreContact.RawPrincipal},
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
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id": contact.UserID,
		}
		rErr := coreerrors.NewNoUserFoundError(fields, true)
		evtString := fmt.Sprintf("no user found with id: %s", contact.UserID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	span.AddEvent("contact added")
	return nil
}

func (ur userRepo) UpdateContact(ctx context.Context, contact *models.Contact, modifiedByID string) errors.RichError {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "UpdateContact", ur.GetType())
	defer span.End()
	contact.AuditData.ModifiedByID = nullable.NullableString{}
	contact.AuditData.ModifiedByID.Set(modifiedByID)
	contact.AuditData.ModifiedOnDate = nullable.NullableTime{}
	contact.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	contactID, err := primitive.ObjectIDFromHex(contact.ID)
	if err != nil {
		rErr := coreerrors.NewFailedToParseObjectIDError(contact.ID, err, true)
		evtString := fmt.Sprintf("%s contact id: %s", rErr.GetErrorMessage(), contact.ID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	oid, err := primitive.ObjectIDFromHex(contact.UserID)
	if err != nil {
		rErr := coreerrors.NewFailedToParseObjectIDError(contact.UserID, err, true)
		evtString := fmt.Sprintf("%s user id: %s", rErr.GetErrorMessage(), contact.UserID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "contacts.id", Value: contactID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "contacts.$.name", Value: contact.Name.GetPointerCopy()},
			{Key: "contacts.$.principal", Value: contact.Principal},
			{Key: "contacts.$.rawPrincipal", Value: contact.RawPrincipal},
			{Key: "contacts.$.type", Value: contact.Type},
			{Key: "contacts.$.isPrimary", Value: contact.IsPrimary},
			{Key: "contacts.$.confirmedDate", Value: contact.ConfirmedDate.GetPointerCopy()},
			{Key: "contacts.$.modifiedById", Value: contact.AuditData.ModifiedByID.GetPointerCopy()},
			{Key: "contacts.$.modifiedOnDate", Value: contact.AuditData.ModifiedOnDate.GetPointerCopy()},
		}},
	}
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update, nil) //(ctx, contact, nil)
	if err != nil {
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id":        contact.UserID,
			"contact.id": contact.ID,
		}
		rErr := coreerrors.NewNoContactFoundError(fields, true)
		evtString := fmt.Sprintf("no contact found with id: %s", contact.ID)
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return rErr
	}
	span.AddEvent("contact updated")
	return nil
}

func (ur userRepo) GetExistingConfirmedContactsCountByPrincipalAndType(ctx context.Context, contactType, contactPrincipal string) (int64, errors.RichError) {
	span := apptelemetry.CreateRepoFunctionSpan(ctx, ur.GetName(), "GetExistingConfirmedContactsCountByPrincipalAndType", ur.GetType())
	defer span.End()
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
	numConfirmedContacts, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).CountDocuments(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.AddEvent("no contacts found")
			return 0, nil
		}
		rErr := coreerrors.NewRepoQueryFailedError(err, true)
		evtString := fmt.Sprintf("repo query failed: %s", rErr.GetErrors()[0].Error())
		apptelemetry.SetSpanOriginalError(&span, rErr, evtString)
		return 0, rErr
	}
	span.AddEvent("number of confirmed contacts retreived")
	return numConfirmedContacts, nil
}
