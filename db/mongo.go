package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBType struct {
	MongoClient *mongo.Client
}

func MongoDB() *MongoDBType {
	return &MongoDBType{}
}

func (db *MongoDBType) Connect() {
	var err error

	// Get MongoDB URI from environment
	mongoURI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(mongoURI)
	db.MongoClient, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = db.MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func (db *MongoDBType) GetCollection(collection string) *mongo.Collection {
	database := os.Getenv("MONGO_DATABASE")
	return db.MongoClient.Database(database).Collection(collection)
}
