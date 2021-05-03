package database

import (
	"context"
	"fmt"
	"log"
	"os"
	errors "server/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type database struct {
	Database *mongo.Database
}

func (d *database) Find(ctx context.Context, coll string, filter interface{}) (*mongo.Cursor, error) {
	return d.Database.Collection(coll).Find(ctx, filter)
}

func (d *database) InsertOne(ctx context.Context, coll string, docs ...interface{}) (*mongo.InsertManyResult, error) {
	return d.Database.Collection(coll).InsertMany(ctx, docs)
}

func (d *database) FindOne(ctx context.Context, coll string, filter interface{}) *mongo.SingleResult {
	return d.Database.Collection(coll).FindOne(ctx, filter)
}

func (d *database) UpdateOne(ctx context.Context, coll string, filter, update interface{}) (*mongo.UpdateResult, error) {
	return d.Database.Collection(coll).UpdateOne(ctx, filter, update)
}

func (d *database) FindOneAndUpdate(ctx context.Context, coll string, filter, update interface{}) *mongo.SingleResult {
	return d.Database.Collection(coll).FindOneAndUpdate(ctx, filter, update)
}

var Database database
var Bucket *gridfs.Bucket

func Connect() {
	connectionString := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(errors.Wrap(err, err.Error()))
	}

	fmt.Println("Connected to MongoDB!")

	Database = database{client.Database(dbName)}

	fmt.Println("Collection instance created!")

	// Initalizing the gridFS
    Bucket, err := gridfs.NewBucket(
        client.Database("myfiles"),
    )
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

	fmt.Println("GridFS bucket created!")
	fmt.Print(Bucket)
}
