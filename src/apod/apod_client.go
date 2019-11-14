package apod

import (
	"io/ioutil"
	"log"
	"net/http"
)

// set up vault to store API key
func GetSinglePicture() {
	httpResp, err := http.Get("https://api.nasa.gov/planetary/apod?api_key=@@@")
	if err != nil {
		log.Fatal(err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
