package mongodb

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoDB *mongo.Database
var m sync.Mutex

func InitMongoClient(uri, dbName string) error {
	m.Lock()
	defer m.Unlock()

	if mongoClient != nil {
		return nil
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
		return err
	}

	mongoClient = client
	mongoDB = client.Database(dbName)
	return nil
}

func GetMongoDatabase() *mongo.Database {
	if mongoDB == nil {
		panic("MongoDB not initialized. Call InitMongoClient first.")
	}
	return mongoDB
}

func CloseMongoClient() {
	m.Lock()
	defer m.Unlock()

	if mongoClient != nil {
		_ = mongoClient.Disconnect(context.TODO())
		mongoClient = nil
		mongoDB = nil
	}
}
