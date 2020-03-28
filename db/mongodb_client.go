package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var databaseAPOD *mongo.Database

// Collection objects that are used within the package
var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection
var savesCollection *mongo.Collection

const dbName = "apodDB"

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

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

	usersCollection = databaseAPOD.Collection("users")
	sessionsCollection = databaseAPOD.Collection("sessions")
	savesCollection = databaseAPOD.Collection("saves")
}
