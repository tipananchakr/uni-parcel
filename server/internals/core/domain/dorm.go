package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dorm struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code      string             `json:"code" bson:"code"`
	Name      string             `json:"name" bson:"name"`
	IsDeleted bool               `json:"is_deleted" bson:"is_deleted"`
}

type DormUpdate struct {
	Code *string
	Name *string
}
