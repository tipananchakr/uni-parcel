package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Role     string             `bson:"role" json:"role"`
}
