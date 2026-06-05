package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates required MongoDB indexes.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	objects := db.Collection("objects")
	events := db.Collection("object_events")

	// Index on owner_user_id
	_, err := objects.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"owner_user_id": 1},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		return err
	}

	// Index on tags
	_, err = objects.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"tags": 1},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		return err
	}

	// Index on object_id for events
	_, err = events.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"object_id": 1},
		Options: options.Index().SetBackground(true),
	})
	if err != nil {
		return err
	}

	log.Println("MongoDB indexes ensured")
	return nil
}
