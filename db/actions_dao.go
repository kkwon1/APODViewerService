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
func NewUserActionDAO(savesCollection *mongo.Collection, likesCollection *mongo.Collection) ActionsDAO {
	return &actionsDAO{
		savesCollection: savesCollection,
		likesCollection: likesCollection,
	}
}

// SaveApod adds a new record of a save action in the database
func (dao *actionsDAO) SaveApod(ctx context.Context, userAction *models.UserAction) error {
	if !saveApodExists(ctx, userAction, dao.savesCollection) {
		_, err := dao.savesCollection.InsertOne(ctx, userAction)
		return err
	}
	return nil
}

// LikeApod adds a new record of a like action in the database
func (dao *actionsDAO) LikeApod(ctx context.Context, userAction *models.UserAction) error {
	if !likeApodExists(ctx, userAction, dao.likesCollection) {
		_, err := dao.likesCollection.InsertOne(ctx, userAction)
		return err
	}
	return nil
}

//TODO: In the future, support multiple saved image deletion
// UnsaveApod removes the save action record with given date
func (dao *actionsDAO) UnsaveApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := dao.savesCollection.DeleteOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	return err
}

// UnlikeApod removes the like action record with given date
func (dao *actionsDAO) UnlikeApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := dao.likesCollection.DeleteOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	return err
}

func saveApodExists(ctx context.Context, userAction *models.UserAction, savesCollection *mongo.Collection) bool {
	singleResult := savesCollection.FindOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	if singleResult.Err() != nil {
		return false
	}

	return true
}

func likeApodExists(ctx context.Context, userAction *models.UserAction, likesCollection *mongo.Collection) bool {
	singleResult := likesCollection.FindOne(ctx, bson.M{"apoddate": userAction.ApodDate, "userid": userAction.UserID})
	if singleResult.Err() != nil {
		return false
	}

	return true
}
