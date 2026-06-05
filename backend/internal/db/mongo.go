package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo connects to MongoDB and returns a client.
func ConnectMongo(uri string) (*mongo.Client, error) {
	log.Println("Connecting to MongoDB...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping to ensure connection is alive
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
