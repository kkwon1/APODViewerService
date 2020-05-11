package db

import (
	"context"

	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ActionsDAO is a DAO for interacting with user action data
type ActionsDAO interface {
	SaveApod(context.Context, *models.UserAction) error
	LikeApod(context.Context, *models.UserAction) error
	UnsaveApod(context.Context, *models.UserAction) error
	UnlikeApod(context.Context, *models.UserAction) error
}

type actionsDAO struct {
	savesCollection *mongo.Collection
	likesCollection *mongo.Collection
}

// NewUserActionDAO returns the actionsDAO object that implements the interface
func NewUserActionDAO() ActionsDAO {
	return &actionsDAO{
		savesCollection: savesCollection,
		likesCollection: likesCollection,
	}
}

// SaveApod adds a new record of a save action in the database
func (u *actionsDAO) SaveApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := savesCollection.InsertOne(ctx, userAction)
	return err
}

// LikeApod adds a new record of a like action in the database
func (u *actionsDAO) LikeApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := likesCollection.InsertOne(ctx, userAction)
	return err
}

//TODO: In the future, support multiple saved image deletion
// UnsaveApod removes the save action record with given date
func (u *actionsDAO) UnsaveApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := savesCollection.DeleteOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	return err
}

// UnlikeApod removes the like action record with given date
func (u *actionsDAO) UnlikeApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := likesCollection.DeleteOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	return err
}
