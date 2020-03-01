package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkwon1/APODViewerService/src/db"
	"github.com/kkwon1/APODViewerService/src/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection *mongo.Collection

func init() {
	var client = db.GetClient()
	usersCollection = client.Database("apodDB").Collection("users")
}

// User type containing Username, Email and Password
type User struct {
	Username string
	Email    string
	Password string
}

// CreateUser stores a new user in the DB given a username, email and password.
// TODO: Make sure user name and email are unique before inserting into DB
// TODO: Re-factor some code in this file to some helper class
func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u User
	decodeError := json.NewDecoder(r.Body).Decode(&u)
	if decodeError != nil {
		log.Fatal(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	bPassword := []byte(u.Password)
	hash, hashError := bcrypt.GenerateFromPassword(bPassword, bcrypt.MinCost)
	if hashError != nil {
		log.Fatal(hashError)
	}

	_, insertError := usersCollection.InsertOne(ctx, bson.D{
		{Key: "username", Value: u.Username},
		{Key: "email", Value: u.Email},
		{Key: "password", Value: string(hash)},
	})

	if insertError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
		log.Fatal(insertError)
		return
	}

	log.Print("Inserted a new user into users collection")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Successfully created user: %s", u.Username)))
}

// DeleteUser deletes user given the correct credentials
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var userToDelete User
	decodeError := json.NewDecoder(r.Body).Decode(&userToDelete)
	if decodeError != nil {
		log.Fatal(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	var retrievedUser User

	findError := usersCollection.FindOne(ctx, bson.M{"username": userToDelete.Username}).Decode(&retrievedUser)

	if findError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userToDelete.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a password"))
		return
	}

	passwordIsCorrect := comparePassword(retrievedUser.Password, userToDelete.Password)

	if passwordIsCorrect {
		res, deleteError := usersCollection.DeleteOne(ctx, bson.M{"username": userToDelete.Username})

		if deleteError != nil {
			log.Fatal(deleteError)
		}

		if res.DeletedCount == 0 {
			fmt.Println("DeleteOne() document not found: ", res)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Service Error"))
		} else {
			fmt.Println("DeleteOne result:", res)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("Successfully deleted user: %s", userToDelete.Username)))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Username or Password is invalid"))
		return
	}
}

// Login checks the user credentials and returns a token if valid.
func Login(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u User
	decodeError := json.NewDecoder(r.Body).Decode(&u)
	if decodeError != nil {
		log.Fatal(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	var retrievedUser User

	findError := usersCollection.FindOne(ctx, bson.M{"username": u.Username}).Decode(&retrievedUser)

	if findError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginIsValid := comparePassword(retrievedUser.Password, u.Password)

	if loginIsValid {
		tokenString := utils.GetJwt()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(tokenString))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

// TestToken is A test endpoint function to check that the given jwt token is valid.
func TestToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	valid := utils.ValidateJwt(tokenString)
	if valid {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("The token is valid"))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("The token is invalid"))
		return
	}
}

// Helper function to compare the password with hash to check that input password is correct
func comparePassword(hashedPwd string, plainPwd string) bool {
	bytePwd := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
