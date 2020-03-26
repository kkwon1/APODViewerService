package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kkwon1/APODViewerService/api/apod"
	"github.com/kkwon1/APODViewerService/api/users"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/users/", users.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/", users.DeleteUser).Methods(http.MethodDelete)
	api.HandleFunc("/users/login", users.Login).Methods(http.MethodPost)
	api.HandleFunc("/apod/batch/", apod.GetBatchImages).Methods(http.MethodGet)

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
