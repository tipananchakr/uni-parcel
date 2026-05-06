package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Major struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	IsDeleted bool               `bson:"is_deleted" json:"is_deleted"`
}

type MajorUpdate struct {
	Code *string
	Name *string
}
