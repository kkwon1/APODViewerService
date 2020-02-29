package apod

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const layoutISO = "2006-01-02"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetSinglePicture() {
	httpResp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s", os.Getenv("NASA_API_KEY")))
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

func GetBatchImages(count int) {
	today := time.Now()
	endDate := today.Format(layoutISO)
	startDate := today.AddDate(0, 0, (count-1)*-1).Format(layoutISO)
	log.Println(endDate)
	log.Println(startDate)

	httpResp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s&start_date=%s&end_date=%s", os.Getenv("NASA_API_KEY"), startDate, endDate))
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
