package apod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

// GetBatchImages will retrieve multiple images from NASA APOD API
func GetBatchImages(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.Query())
	count, convErr := strconv.Atoi(r.URL.Query().Get("count"))

	if convErr != nil {
		log.Fatal(convErr)
	}

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

	log.Printf("Successfully retrieved %d number of images", count)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
