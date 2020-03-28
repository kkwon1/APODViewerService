package db

import (
	"context"

	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Add a likesCollection

// ActionsDAO is a DAO for interacting with user action data
type ActionsDAO interface {
	SaveApod(context.Context, *models.UserAction) error
}

type actionsDAO struct {
	savesCollection *mongo.Collection
}

// NewUserActionDAO returns the actionsDAO object that implements the interface
func NewUserActionDAO() ActionsDAO {
	return &actionsDAO{
		savesCollection: savesCollection,
	}
}

// SaveApod saves a new record of a user action in the database
func (u *actionsDAO) SaveApod(ctx context.Context, userAction *models.UserAction) error {
	_, err := savesCollection.InsertOne(ctx, userAction)
	return err
}
