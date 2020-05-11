package users

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
)

var actionsDAO = db.NewUserActionDAO()

//TODO: Make a single endpoint for like/save and split off in code? Lots of repeated logic

// SaveContent is an endpoint that allows users to save/favourite an APOD of their choosing.
func UserAction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debugln("Starting UserAction call")
	var userAction *models.UserAction

	decodeError := json.NewDecoder(r.Body).Decode(&userAction)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected payload"))
		log.Errorln(decodeError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	idToken := splitToken[1]

	verified, verify_err := VerifyToken(ctx, idToken)
	if verify_err != nil || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Verification Failed"))
		log.Errorln(verify_err)
		return
	}

	switch userAction.Action {
	case "save":
		insertError := actionsDAO.SaveApod(ctx, userAction)
		if insertError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Service Error"))
			log.Fatal(insertError)
			return
		}
	case "like":
		insertError := actionsDAO.LikeApod(ctx, userAction)
		if insertError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Service Error"))
			log.Fatal(insertError)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid action in request body"))
		return
	}

	log.Printf("Username: %s, Successfully completed user action: %s", userAction.UserID, userAction.Action)
}
