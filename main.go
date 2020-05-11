package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkwon1/APODViewerService/api/apod"
	"github.com/kkwon1/APODViewerService/api/users"
	"github.com/kkwon1/APODViewerService/utils"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HelloServer)
	api := r.PathPrefix("/api/v1").Subrouter()

	tokenVerifier := utils.NewTokenVerifier()
	userAction := users.NewUserAction(tokenVerifier)
	userDataRetriever := users.NewUserDataRetriever(tokenVerifier)

	api.HandleFunc("/users/action/", userAction.ApplyAction).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/apod/batch/", apod.GetBatchImages).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/users/data/", userDataRetriever.RetrieveUserData).Methods(http.MethodPost, http.MethodOptions)

	r.Use(CORS)

	srv := &http.Server{
		Addr:    ":8081",
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
	fmt.Fprintf(w, "Hello World!")
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
