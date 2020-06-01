package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/kkwon1/APODViewerService/api/apod"
	"github.com/kkwon1/APODViewerService/api/users"
	"github.com/kkwon1/APODViewerService/db"
	"github.com/kkwon1/APODViewerService/utils"

	"github.com/gorilla/mux"
)

var userAction users.UserAction
var userDataRetriever users.UserDataRetriever
var apodClient apod.ApodClient

func dependencyInit() {
	apodClient = apod.NewApodClient()
	tokenVerifier := utils.NewTokenVerifier()
	apodDb := db.NewMongoClient().GetApodDB(context.Background(), os.Getenv("MONGODB_URI"))

	// Initialize the collections that will be used in DAOs
	savesCollection := apodDb.Collection("saves")
	likesCollection := apodDb.Collection("likes")
	userActionDAO := db.NewUserActionDAO(savesCollection, likesCollection)
	userDataDAO := db.NewUserDataDAO(savesCollection, likesCollection)

	userAction = users.NewUserAction(tokenVerifier, userActionDAO)
	userDataRetriever = users.NewUserDataRetriever(tokenVerifier, userDataDAO)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	r := mux.NewRouter()
	r.HandleFunc("/", HelloServer)
	api := r.PathPrefix("/api/v1").Subrouter()

	dependencyInit()

	api.HandleFunc("/users/action/", userAction.ApplyAction).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/apod/batch/", apodClient.GetBatchImages).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/users/data/", userDataRetriever.RetrieveUserData).Methods(http.MethodPost, http.MethodOptions)

	r.Use(CORS)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		// e.g. shut down connection to db etc.
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}

// TODO: Remove
// Temporary endpoint for debugging purposes
func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! Heroku Test")
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers:", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PATCH, PUT, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
	})
}
