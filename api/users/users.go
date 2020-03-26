package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kkwon1/APODViewerService/db"
	uuid "github.com/satori/go.uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection *mongo.Collection
var sessionsCollection *mongo.Collection

func init() {
	var client = db.GetClient()
	usersCollection = client.Database("apodDB").Collection("users")
	sessionsCollection = client.Database("apodDB").Collection("sessions")
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

	retrievedUser, findError := checkUserExists(ctx, userToDelete.Username)

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

	if !passwordIsCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Username or Password is invalid"))
		return
	}

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

	retrievedUser, findError := checkUserExists(ctx, u.Username)

	if findError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User does not exist"))
		return
	}

	loginIsValid := comparePassword(retrievedUser.Password, u.Password)

	if !loginIsValid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Username or Password is invalid"))
		return
	}

	// Create new random session token
	sessionToken := uuid.NewV4().String()
	expiryTime := time.Now().Add(30 * time.Minute)

	// Store session to db along with user
	// TODO: First check that user already has a session token, then just update session ID.
	// If not, add new record entirely
	_, insertError := sessionsCollection.InsertOne(ctx, bson.D{
		{Key: "username", Value: u.Username},
		{Key: "sessionToken", Value: sessionToken},
		{Key: "expiryTime", Value: expiryTime},
	})

	if insertError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
		log.Fatal(insertError)
		return
	}

	// Set client cookie with session token
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   sessionToken,
		Expires: expiryTime,
	})
}

// Function to simply check whether user already exists in DB
func checkUserExists(ctx context.Context, username string) (User, error) {
	var retrievedUser User

	findError := usersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&retrievedUser)

	return retrievedUser, findError
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
