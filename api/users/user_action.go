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
	"github.com/kkwon1/APODViewerService/utils"
)

type UserAction interface {
	ApplyAction(w http.ResponseWriter, r *http.Request)
}

type userAction struct {
	tokenVerifier utils.TokenVerifier
	actionsDAO    db.ActionsDAO
}

func NewUserAction(tokenVerifier utils.TokenVerifier, actionsDAO db.ActionsDAO) UserAction {
	return &userAction{
		tokenVerifier: tokenVerifier,
		actionsDAO:    actionsDAO,
	}
}

// ApplyAction is an endpoint that allows users to like or save an APOD
func (ua *userAction) ApplyAction(w http.ResponseWriter, r *http.Request) {
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

	verified, verify_err := ua.tokenVerifier.VerifyToken(ctx, idToken)
	if verify_err != nil || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Verification Failed"))
		log.Errorln(verify_err)
		return
	}

	var err error
	switch userAction.Action {
	case "save":
		err = ua.actionsDAO.SaveApod(ctx, userAction)
		log.Println("Successfully saved APOD")
	case "unsave":
		err = ua.actionsDAO.UnsaveApod(ctx, userAction)
		log.Println("Successfully unsaved APOD")
	case "like":
		err = ua.actionsDAO.LikeApod(ctx, userAction)
		log.Println("Successfully liked APOD")
	case "unlike":
		err = ua.actionsDAO.UnlikeApod(ctx, userAction)
		log.Println("Successfully unliked APOD")
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid action in request body"))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Service Error"))
		log.Errorln(err)
		return
	} else if err == nil {
		apodData := convertUserActionToApodObject(userAction)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(apodData)
	}

	log.Printf("UserID: %s, Successfully completed user action: %s", userAction.UserID, userAction.Action)
}

func convertUserActionToApodObject(userAction *models.UserAction) *models.ApodObject {
	apodObject := &models.ApodObject{
		UserID:      userAction.UserID,
		ApodURL:     userAction.ApodURL,
		ApodName:    userAction.ApodName,
		ApodDate:    userAction.ApodDate,
		MediaType:   userAction.MediaType,
		Description: userAction.Description,
		ActionDate:  userAction.ActionDate,
	}

	return apodObject
}
