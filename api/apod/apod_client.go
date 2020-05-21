package apod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const layoutISO = "2006-01-02"

// GetBatchImages will retrieve multiple images from NASA APOD API
func GetBatchImages(w http.ResponseWriter, r *http.Request) {
	count, convErr := strconv.Atoi(r.URL.Query().Get("count"))

	if convErr != nil {
		log.Errorln(convErr)
	}

	page, convErr := strconv.Atoi(r.URL.Query().Get("page"))

	if convErr != nil {
		page = 0
	}

	startDate, endDate := getStartEndDates(count, page)
	log.Printf("End date: %s", endDate)
	log.Printf("Start date: %s", startDate)

	httpResp, err := http.Get(fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s&start_date=%s&end_date=%s", os.Getenv("NASA_API_KEY"), startDate, endDate))
	if err != nil {
		log.Errorln(err)
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorln(err)
	}

	log.Printf("Successfully retrieved %d number of images", count)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func getStartEndDates(count int, page int) (string, string) {
	// page * count of batch images will be the offset of date
	offset := page * count
	// Get today's date in UTC (That is how APOD operates)
	// convert to specified ISO, and find start date to call API and retrieve
	// the correct number of images
	today := time.Now()
	endDateTime := today.AddDate(0, 0, offset*-1)
	endDate := endDateTime.Format(layoutISO)
	startDate := endDateTime.AddDate(0, 0, (count-1)*-1).Format(layoutISO)

	return startDate, endDate
}
