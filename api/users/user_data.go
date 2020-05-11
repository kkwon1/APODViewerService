package users

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
	"github.com/kkwon1/APODViewerService/utils"
	log "github.com/sirupsen/logrus"
)

var userDataDAO = db.NewUserDataDAO()

type UserDataRetriever interface {
	RetrieveUserData(w http.ResponseWriter, r *http.Request)
}

type userDataRetriever struct {
	tokenVerifier utils.TokenVerifier
}

func NewUserDataRetriever(tokenVerifier utils.TokenVerifier) UserDataRetriever {
	return &userDataRetriever{
		tokenVerifier: tokenVerifier,
	}
}

func (udr *userDataRetriever) RetrieveUserData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userDataRetrievalModel *models.UserAction

	decodeError := json.NewDecoder(r.Body).Decode(&userDataRetrievalModel)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unexpected payload"))
		log.Errorln(decodeError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	idToken := splitToken[1]

	verified, verifyError := udr.tokenVerifier.VerifyToken(ctx, idToken)
	if verifyError != nil || !verified {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Verification Failed"))
		log.Errorln(verifyError)
		return
	}

	userSaves, getSaveError := userDataDAO.GetUserSaves(ctx, userDataRetrievalModel.UserID)
	if getSaveError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Verification Failed"))
		log.Errorln(verifyError)
		return
	}

	userLikes, _ := userDataDAO.GetUserLikes(ctx, userDataRetrievalModel.UserID)

	userData := models.UserData{
		UserSaves: userSaves,
		UserLikes: userLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
