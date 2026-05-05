package mongodb

import (
	"context"
	"errors"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DormRepository struct {
	collection *mongo.Collection
}

func NewDormRepository(collection *mongo.Collection) *DormRepository {
	return &DormRepository{
		collection: collection,
	}
}

func (d *DormRepository) GetAll(ctx context.Context) ([]*domain.Dorm, error) {
	cursor, err := d.collection.Find(ctx, bson.M{"is_deleted": false})

	if err != nil {
		return nil, errors.New("failed to fetch dorms: " + err.Error())
	}

	defer cursor.Close(ctx)

	var dorms []*domain.Dorm
	if err := cursor.All(ctx, &dorms); err != nil {
		return nil, errors.New("failed to decode dorms: " + err.Error())
	}

	return dorms, nil
}

func (d *DormRepository) GetByID(ctx context.Context, id string) (*domain.Dorm, error) {
	dormId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid dorm ID: " + err.Error())
	}

	var dorm domain.Dorm
	err = d.collection.FindOne(ctx, bson.M{
		"_id":        dormId,
		"is_deleted": bson.M{"$ne": true},
	}).Decode(&dorm)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("dorm not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch dorm: " + err.Error())
	}

	return &dorm, nil
}

func (d *DormRepository) Create(ctx context.Context, dorm *domain.Dorm) error {
	dormDocument := domain.Dorm{
		Code:      dorm.Code,
		Name:      dorm.Name,
		IsDeleted: false,
	}

	result, err := d.collection.InsertOne(ctx, dormDocument)
	if err != nil {
		return errors.New("failed to create dorm: " + err.Error())
	}

	dorm.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (d *DormRepository) Update(ctx context.Context, id string, update domain.DormUpdate) error {
	dormId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid dorm ID: " + err.Error())
	}

	set := bson.M{}
	if update.Code != nil {
		set["code"] = *update.Code
	}
	if update.Name != nil {
		set["name"] = *update.Name
	}

	if len(set) == 0 {
		return errors.New("no fields to update")
	}

	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": dormId, "is_deleted": bson.M{"$ne": true}},
		bson.M{"$set": set},
	)
	if err != nil {
		return errors.New("failed to update dorm: " + err.Error())
	}

	if result.MatchedCount == 0 {
		return errors.New("dorm not found")
	}

	return nil
}

func (d *DormRepository) Delete(ctx context.Context, id string) error {
	dormId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid dorm ID: " + err.Error())
	}

	_, err = d.collection.UpdateOne(ctx, bson.M{"_id": dormId, "is_deleted": bson.M{"$ne": true}}, bson.M{
		"$set": bson.M{"is_deleted": true},
	})
	if err != nil {
		return errors.New("failed to delete dorm: " + err.Error())
	}

	return nil
}
