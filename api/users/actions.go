package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/bson"
)

var actionsDAO = db.NewUserActionDAO()

// SaveContent is an endpoint that allows users to save/favourite an APOD of their choosing.
func SaveContent(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var userAction *models.UserAction
	decodeError := json.NewDecoder(r.Body).Decode(&userAction)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected payload"))
		log.Fatal(decodeError)
		return
	}

	// Read the session token cookie from request
	cookie, cookieError := r.Cookie("sessionToken")
	if cookieError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session token is missing. Please login"))
		log.Fatal(cookieError)
	}

	// Parse string value
	sessionToken := cookie.Value

	log.Println(userAction)
	if !sessionIsValid(ctx, sessionToken, userAction.Username) {
		// TODO: Add some kind of refresh mechanism for session token
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Session token has expired. Please login again"))
		log.Fatal("Session token has expired.")
		return
	}

	insertError := actionsDAO.SaveApod(ctx, userAction)

	if insertError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
		log.Fatal(insertError)
		return
	}
}

// TODO: Implement this function
/*
func LikeContent()
*/

func sessionIsValid(ctx context.Context, sessionToken string, username string) bool {
	session, err := sessionsDAO.FindOne(ctx, bson.M{"username": username})
	if err != nil {
		log.Print(err)
	}

	currentTime := time.Now().Unix()

	log.Println(session)
	log.Println(currentTime)
	log.Println(session.ExpiryTime)
	log.Println(session.SessionToken)
	log.Println(sessionToken)
	return session.ExpiryTime > currentTime && (session.SessionToken == sessionToken)
}
