package db

import (
	"context"

	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UsersDAO is an interface for accessing the database
type UsersDAO interface {
	FindOne(context.Context, interface{}) (*models.User, error)
	CreateUser(context.Context, *models.User) error
	DeleteByUsername(context.Context, string) (*mongo.DeleteResult, error)
}

type usersDao struct {
	usersCollection *mongo.Collection
}

// NewUserDAO returns the usersDao object that implements the interface
func NewUserDAO() UsersDAO {
	return &usersDao{
		usersCollection: usersCollection,
	}
}

// FindOne returns a user if found in db. If not, returns an error
func (u *usersDao) FindOne(ctx context.Context, filter interface{}) (*models.User, error) {
	user := &models.User{}
	err := usersCollection.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateUser creates a user in the database
func (u *usersDao) CreateUser(ctx context.Context, user *models.User) error {
	_, err := usersCollection.InsertOne(ctx, user)
	return err
}

// DeleteByUsername deletes a user in the db given a username.
// TODO: This will probably be a cascading deletion operation in the future.
func (u *usersDao) DeleteByUsername(ctx context.Context, username string) (*mongo.DeleteResult, error) {
	res, err := usersCollection.DeleteOne(ctx, bson.M{"username": username})
	return res, err
}
