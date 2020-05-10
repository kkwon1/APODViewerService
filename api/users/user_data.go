package users

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/models"
)

var userDataDAO = db.NewUserDataDAO()

func RetrieveUserData(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	userSaves, _ := userDataDAO.GetUserSaves(ctx, "iVKyHd6Rs1ParGxKwHBDTABMEBv1")
	userLikes, _ := userDataDAO.GetUserLikes(ctx, "iVKyHd6Rs1ParGxKwHBDTABMEBv1")

	userData := models.UserData{
		UserSaves: userSaves,
		UserLikes: userLikes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

func getUserLikes() {

}

func getUserSaves() {

}

// TODO: Stub. Once settings implemented, this should return the saved settings.
func getUserProfileSettings() {

}
