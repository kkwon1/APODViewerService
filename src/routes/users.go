package routes

import (
	"APODViewerService/src/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var UsersCollection *mongo.Collection

func init() {
	var client = db.GetClient()
	UsersCollection = client.Database("apodDB").Collection("users")
}

type User struct {
	Username string
	Email    string
}

type UsernameStruct struct {
	Username string
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u User
	decodeError := json.NewDecoder(r.Body).Decode(&u)
	if decodeError != nil {
		log.Fatal(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	_, insertError := UsersCollection.InsertOne(ctx, bson.D{
		{Key: "username", Value: u.Username},
		{Key: "email", Value: u.Email},
	})

	if insertError != nil {
		log.Fatal(insertError)
	}

	log.Print("Inserted a new user into users collection")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var usernameToDelete UsernameStruct
	decodeError := json.NewDecoder(r.Body).Decode(&usernameToDelete)
	if decodeError != nil {
		log.Fatal(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	res, deleteError := UsersCollection.DeleteOne(ctx, bson.M{"username": usernameToDelete.Username})

	if deleteError != nil {
		log.Fatal(deleteError)
	}

	if res.DeletedCount == 0 {
		fmt.Println("DeleteOne() document not found: ", res)
	} else {
		fmt.Println("DeleteOne result:", res)
	}
}
