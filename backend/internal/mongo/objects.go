package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Object represents a digital object stored in MongoDB.
type Object struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerUserID string             `bson:"owner_user_id" json:"owner_user_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Tags        []string           `bson:"tags" json:"tags"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ObjectEvent stores history of changes to an object.
type ObjectEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ObjectID  primitive.ObjectID `bson:"object_id" json:"object_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Type      string             `bson:"type" json:"type"` // created, updated, transferred, deleted
	Payload   interface{}        `bson:"payload" json:"payload"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
