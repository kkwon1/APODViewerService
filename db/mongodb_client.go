package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseAPOD *mongo.Database

// Collection objects that are used within the package
var savesCollection *mongo.Collection
var likesCollection *mongo.Collection

const dbName = "apodDB"

func init() {
	// Set client options
	mongo_uri := fmt.Sprintf(os.Getenv("MONGODB_URI"))
	clientOptions := options.Client().ApplyURI(mongo_uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	databaseAPOD = client.Database(dbName)

	savesCollection = databaseAPOD.Collection("saves")
	likesCollection = databaseAPOD.Collection("likes")
}
