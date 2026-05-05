package mongodb

import (
	"context"

	"github.com/tipananchakr/uni-parcel/internals/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	document := domain.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
	}

	result, err := r.collection.InsertOne(ctx, document)
	if mongo.IsDuplicateKeyError(err) {
		return domain.User{}, domain.ErrEmailAlreadyExists
	}
	if err != nil {
		return domain.User{}, domain.ErrInternalServerError
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return domain.User{}, domain.ErrInternalServerError
	}

	document.ID = insertedID

	return document, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.FindOne(ctx, bson.M{"email": email})
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, domain.ErrInvalidTodoID
	}

	return r.FindOne(ctx, bson.M{"_id": objectID})
}

func (r *UserRepository) FindOne(ctx context.Context, filter bson.M) (domain.User, error) {
	var document domain.User
	err := r.collection.FindOne(ctx, filter).Decode(&document)
	if err == mongo.ErrNoDocuments {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, domain.ErrInternalServerError
	}

	return document, nil
}
