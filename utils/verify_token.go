package utils

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"

	"google.golang.org/api/option"
)

type TokenVerifier interface {
	VerifyToken(ctx context.Context, idToken string) (bool, error)
}

type tokenVerifier struct {
}

func NewTokenVerifier() TokenVerifier {
	return &tokenVerifier{}
}

func (tv *tokenVerifier) VerifyToken(ctx context.Context, idToken string) (bool, error) {
	credentials_file := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	opt := option.WithCredentialsFile(credentials_file)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Errorf("error initializing app: %v", err)
		return false, fmt.Errorf("error initializing app: %v", err)
	}

	client, auth_err := app.Auth(ctx)
	if auth_err != nil {
		log.Errorf("error getting Auth client: %v\n", auth_err)
		return false, fmt.Errorf("error getting Auth client: %v\n", auth_err)
	}

	token, verify_err := client.VerifyIDToken(ctx, idToken)
	if verify_err != nil {
		log.Errorf("error verifying ID token: %v\n", verify_err)
		return false, fmt.Errorf("error verifying ID token: %v\n", verify_err)
	}

	log.Debugf("Verified ID token: %v\n", token)

	return true, nil
}
