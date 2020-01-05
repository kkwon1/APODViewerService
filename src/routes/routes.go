package routes

import (
	"APODViewerService/src/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	dbClient := db.GetClient()
	apodDatabase := dbClient.Database("apodDB")
	usersCollection := apodDatabase.Collection("users")

	_, err := usersCollection.InsertOne(ctx, bson.D{
		{Key: "username", Value: "kevintest"},
		{Key: "email", Value: "test@test.com"},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Inserted a new user into users collection")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	dbClient := db.GetClient()
	apodDatabase := dbClient.Database("apodDB")
	usersCollection := apodDatabase.Collection("users")

	res, err := usersCollection.DeleteOne(ctx, bson.M{"username": "kevintest"})

	if err != nil {
		log.Fatal(err)
	}

	if res.DeletedCount == 0 {
		fmt.Println("DeleteOne() document not found: ", res)
	} else {
		fmt.Println("DeleteOne result:", res)
	}
}
