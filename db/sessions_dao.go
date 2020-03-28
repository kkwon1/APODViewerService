package db

import (
	"context"

	"github.com/kkwon1/APODViewerService/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SessionsDAO is an interface for accessing the database
type SessionsDAO interface {
	FindOne(context.Context, interface{}) (*models.Session, error)
	CreateSessionRecord(context.Context, *models.Session) error
	UpdateSessionToken(context.Context, string, string) error
}

type sessionsDao struct {
	sessionsCollection *mongo.Collection
}

// NewSessionsDAO returns the sessionsDao object that implements the interface
func NewSessionsDAO() SessionsDAO {
	return &sessionsDao{
		sessionsCollection: sessionsCollection,
	}
}

// FindOne returns a session record if found in db
func (s *sessionsDao) FindOne(ctx context.Context, filter interface{}) (*models.Session, error) {
	session := &models.Session{}
	err := sessionsCollection.FindOne(ctx, filter).Decode(session)

	return session, err
}

// CreateSessionRecord a session record in the database for a given user
func (s *sessionsDao) CreateSessionRecord(ctx context.Context, session *models.Session) error {
	_, err := sessionsCollection.InsertOne(ctx, session)
	return err
}

// UpdateSessionToken updates the session token for a given user
func (s *sessionsDao) UpdateSessionToken(ctx context.Context, username string, sessionToken string) error {
	_, err := sessionsCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"sessionToken": sessionToken})
	return err
}
