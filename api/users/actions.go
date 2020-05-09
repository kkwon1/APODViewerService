package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
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
		log.Errorln(decodeError)
		return
	}

	verified, verify_err := VerifyToken(ctx, userAction.IDToken)
	if verify_err != nil || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Verification Failed"))
		log.Errorln(verify_err)
		return
	}

	log.Println("Token verified!")

	log.Println(userAction)

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
