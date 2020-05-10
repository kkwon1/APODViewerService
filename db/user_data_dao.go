package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ActionsDAO is a DAO for interacting with user action data
type UserDataDAO interface {
	GetUserSaves(context.Context, string) ([]*models.ApodObject, error)
	GetUserLikes(context.Context, string) ([]*models.ApodObject, error)
}

type userDataDAO struct {
	savesCollection *mongo.Collection
	likesCollection *mongo.Collection
}

// NewUserDataDAO returns the userDataDAO object that implements the interface
func NewUserDataDAO() UserDataDAO {
	return &userDataDAO{
		savesCollection: savesCollection,
		likesCollection: likesCollection,
	}
}

// SaveApod adds a new record of a save action in the database
func (u *userDataDAO) GetUserSaves(ctx context.Context, userID string) ([]*models.ApodObject, error) {
	cursor, err := savesCollection.Find(ctx, bson.M{"userid": userID})
	var results []*models.ApodObject

	// iterate through all documents
	for cursor.Next(ctx) {
		var object *models.ApodObject

		// Decode the document
		if decode_error := cursor.Decode(&object); decode_error != nil {
			log.Println("cursor.Decode ERROR:", decode_error)
		}
		results = append(results, object)
	}
	return results, err
}

// LikeApod adds a new record of a like action in the database
func (u *userDataDAO) GetUserLikes(ctx context.Context, userID string) ([]*models.ApodObject, error) {
	cursor, err := likesCollection.Find(ctx, bson.M{"userid": userID})
	var results []*models.ApodObject

	// iterate through all documents
	for cursor.Next(ctx) {
		var object *models.ApodObject

		// Decode the document
		if decode_error := cursor.Decode(&object); decode_error != nil {
			log.Println("cursor.Decode ERROR:", decode_error)
		}
		results = append(results, object)
	}
	return results, err
}
