package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"log"
	"time"
)

type MongoDatabase struct {
	DbUrl  string
	DbName string
	client *mongo.Client
	Logger *zap.SugaredLogger
}

func (db *MongoDatabase) GetClient() *mongo.Client {
	return db.client
}

func (db *MongoDatabase) Connect() error {
	if db.DbUrl == "" {
		db.Logger.Errorw("DbUrl is empty")
		return errors.New("DbUrl is empty")
	}
	if db.DbName == "" {
		db.Logger.Errorw("DbName is empty")
		return errors.New("DbName is empty")
	}

	db.Logger.Infow("Connecting to MongoDB...", "url", db.DbUrl)
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the MongoDB client options
	clientOptions := options.Client().ApplyURI(db.DbUrl)

	// Connect to MongoDB
	var err error
	db.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server to verify the connection
	err = db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (db *MongoDatabase) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	db.Logger.Infow("Disconnected from MongoDB.")
}

func (db *MongoDatabase) Find(collection string, filter interface{}, result interface{}) error {
	// Get a collection object from the MongoDB database
	coll := db.client.Database(db.DbName).Collection(collection)

	// Use the Find method of the collection object to execute the query
	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	// Decode the results into the result interface
	err = cursor.All(context.Background(), result)
	if err != nil {
		return err
	}

	return nil
}
