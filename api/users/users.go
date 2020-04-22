package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var usersDAO = db.NewUserDAO()
var sessionsDAO = db.NewSessionsDAO()

// CreateUser stores a new user in the DB given a username, email and password.
// TODO: Make sure user name and email are unique before inserting into DB
// TODO: Re-factor some code in this file to some helper class
func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u *models.User
	decodeError := json.NewDecoder(r.Body).Decode(&u)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected request payload"))
		log.Errorln(decodeError)
		return
	}

	bPassword := []byte(u.Password)
	hash, hashError := bcrypt.GenerateFromPassword(bPassword, bcrypt.MinCost)
	if hashError != nil {
		log.Errorln(hashError)
	}

	u.Password = string(hash)
	insertError := usersDAO.CreateUser(ctx, u)

	if insertError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
		log.Errorln(insertError)
		return
	}

	log.Println("Inserted a new user into users collection")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Successfully created user: %s", u.Username)))
}

// DeleteUser deletes user given the correct credentials
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u *models.User
	decodeError := json.NewDecoder(r.Body).Decode(&u)
	if decodeError != nil {
		log.Errorln(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
		return
	}

	retrievedUser, findError := usersDAO.FindOne(ctx, bson.M{"username": u.Username})

	if findError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if u.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide a password"))
		return
	}

	passwordIsCorrect := comparePassword(retrievedUser.Password, u.Password)

	if !passwordIsCorrect {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Username or Password is invalid"))
		return
	}

	res, deleteError := usersDAO.DeleteByUsername(ctx, u.Username)

	if deleteError != nil {
		log.Errorln(deleteError)
	}

	if res.DeletedCount == 0 {
		fmt.Println("DeleteOne() document not found: ", res)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
	} else {
		fmt.Println("DeleteOne result:", res)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf("Successfully deleted user: %s", u.Username)))
	}
}

// Login checks the user credentials and returns a token if valid.
func Login(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var u *models.User
	decodeError := json.NewDecoder(r.Body).Decode(&u)

	if decodeError != nil {
		log.Errorln(decodeError)
		http.Error(w, decodeError.Error(), http.StatusBadRequest)
	}

	retrievedUser, findError := usersDAO.FindOne(ctx, bson.M{"username": u.Username})

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

	session := &models.Session{
		Username:     u.Username,
		SessionToken: sessionToken,
		ExpiryTime:   expiryTime.Unix(),
	}

	// Store session to db along with user
	// TODO: First check that user already has a session token, then just update session ID.
	// If not, add new record entirely
	userSessionExists := userSessionExists(ctx, u.Username)
	if userSessionExists {
		updateError := sessionsDAO.UpdateSessionToken(ctx, session)

		if updateError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Service Error"))
			log.Errorln(updateError)
			return
		}
	} else {
		insertError := sessionsDAO.CreateSessionRecord(ctx, session)

		if insertError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Service Error"))
			log.Errorln(insertError)
			return
		}
	}

	// Set client cookie with session token
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   sessionToken,
		Expires: expiryTime,
	})
}

func userSessionExists(ctx context.Context, username string) bool {
	_, err := sessionsDAO.FindOne(ctx, bson.M{"username": username})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
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
