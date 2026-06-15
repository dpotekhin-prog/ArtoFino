package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Object represents a digital art asset stored in MongoDB with live economic parameters.
type Object struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerUserID     string             `bson:"owner_user_id" json:"ownerUserId"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Tags            []string           `bson:"tags" json:"tags"`
	BasePrice       float64            `bson:"base_price" json:"basePrice"`
	Currency        string             `bson:"currency" json:"currency"`
	DailyGrowthRate float64            `bson:"daily_growth_rate" json:"dailyGrowthRate"`
	CurrentHolderID string             `bson:"current_holder_id" json:"currentHolderId"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ObjectEvent stores history of changes to an object.
type ObjectEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ObjectID  primitive.ObjectID `bson:"object_id" json:"objectId"`
	UserID    string             `bson:"user_id" json:"userId"`
	Type      string             `bson:"type" json:"type"` // created, updated, transferred, deleted
	Payload   interface{}        `bson:"payload" json:"payload"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// ObjectsRepository handles database operations for the objects collection.
type ObjectsRepository struct {
	collection *mongo.Collection
}

// NewObjectsRepository initializes the repository targeting your "objects" collection.
func NewObjectsRepository(db *mongo.Database) *ObjectsRepository {
	return &ObjectsRepository{
		collection: db.Collection("objects"),
	}
}

// FindByID retrieves a single object document from MongoDB using hex ID string.
func (r *ObjectsRepository) FindByID(ctx context.Context, idStr string) (*Object, error) {
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, err
	}

	var obj Object
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func (r *ObjectsRepository) Create(ctx context.Context, obj *Object) error {
	if obj.ID.IsZero() {
		obj.ID = primitive.NewObjectID()
	}
	if obj.CreatedAt.IsZero() {
		obj.CreatedAt = time.Now()
	}
	obj.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, obj)
	return err
}
