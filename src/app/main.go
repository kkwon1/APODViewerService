package main

import (
	"APODViewerService/src/apod"
	"APODViewerService/src/db"
	"fmt"
	"net/http"
)

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	apod.GetSinglePicture()
	db.ConnectClient()
	http.HandleFunc("/", helloServer)
	http.ListenAndServe(":8081", nil)
}
