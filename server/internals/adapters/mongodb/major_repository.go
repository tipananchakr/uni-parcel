package mongodb

import (
	"context"
	"errors"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MajorRepository struct {
	collection *mongo.Collection
}

func NewMajorRepository(collection *mongo.Collection) *MajorRepository {
	return &MajorRepository{
		collection: collection,
	}
}

func (m *MajorRepository) GetAll(ctx context.Context) ([]*domain.Major, error) {
	cursor, err := m.collection.Find(ctx, bson.M{
		"is_deleted": false,
	})

	if err != nil {
		return nil, errors.New("failed to fetch majors" + err.Error())
	}

	var majors []*domain.Major
	if err := cursor.All(ctx, &majors); err != nil {
		return nil, errors.New("failed to decode majors: " + err.Error())
	}

	return majors, nil
}

func (m *MajorRepository) GetByID(ctx context.Context, id string) (*domain.Major, error) {
	majorId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid major ID: " + err.Error())
	}

	var major domain.Major
	err = m.collection.FindOne(ctx, bson.M{
		"_id": majorId,
		"is_deleted": bson.M{
			"$ne": true,
		},
	}).Decode(&major)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("major not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch major: " + err.Error())
	}

	return &major, nil
}

func (m *MajorRepository) Create(ctx context.Context, major *domain.Major) error {
	majorDocument := domain.Major{
		Code:      major.Code,
		Name:      major.Name,
		IsDeleted: false,
	}

	result, err := m.collection.InsertOne(ctx, majorDocument)
	if err != nil {
		return errors.New("failed to create major: " + err.Error())
	}

	major.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *MajorRepository) Update(ctx context.Context, id string, update domain.MajorUpdate) error {
	majorId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid major ID: " + err.Error())
	}

	updateDoc := bson.M{}
	if update.Code != nil {
		updateDoc["code"] = *update.Code
	}
	if update.Name != nil {
		updateDoc["name"] = *update.Name
	}

	if len(updateDoc) == 0 {
		return errors.New("no fields to update")
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{
		"_id": majorId,
	}, bson.M{
		"$set": updateDoc,
	})

	if err != nil {
		return errors.New("failed to update major: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.New("major not found")
	}

	return nil
}

func (m *MajorRepository) Delete(ctx context.Context, id string) error {
	majorId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid major ID: " + err.Error())
	}

	_, err = m.collection.UpdateOne(ctx, bson.M{
		"_id": majorId,
	}, bson.M{
		"$set": bson.M{
			"is_deleted": true,
		},
	})

	if err != nil {
		return errors.New("failed to delete major: " + err.Error())
	}

	return nil
}
