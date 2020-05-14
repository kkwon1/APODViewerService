package db

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbClient interface {
	GetApodDB(ctx context.Context, mongo_uri string) *mongo.Database
}

type mongoDbClient struct {
}

func NewMongoClient() MongoDbClient {
	return &mongoDbClient{}
}

func (m *mongoDbClient) GetApodDB(ctx context.Context, mongo_uri string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(mongo_uri)
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	log.Printf("Successfully connected to MongoDB at %s, using database %s", mongo_uri, dbName)
	databaseAPOD := client.Database(dbName)

	return databaseAPOD
}
