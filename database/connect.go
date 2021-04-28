package database

import (
	"context"
	"fmt"
	"log"
	"os"
	errors "server/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collection object/instance - used across database package
var Collection *mongo.Collection

// connect to mongodb
func Connect() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collName := os.Getenv("DB_COLLECTION_NAME")

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}

	fmt.Println("Connected to MongoDB!")

	Collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}
