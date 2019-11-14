package main

import (
	"fmt"
	"net/http"
)

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", helloServer)
	http.ListenAndServe(":8081", nil)
}
